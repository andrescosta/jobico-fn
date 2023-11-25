package listener

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Controller struct {
	queueHost string
	events    map[string]*Event
}

type Event struct {
	event  *pb.EventDef
	schema *jsonschema.Schema
}

func New(ctx context.Context) (*Controller, error) {
	client := remote.NewControlClient(env.GetAsString("ctl.host", ""))
	repoClient := remote.NewRepoClient(env.GetAsString("repo.host", ""))

	events := make(map[string]*Event)
	con := Controller{
		events:    events,
		queueHost: env.GetAsString("queue.host", ""),
	}

	pkgs, err := client.GetAllPackages(ctx)
	if err != nil {
		return nil, err
	}
	for _, ps := range pkgs {
		m := ps.MerchantId
		for _, e := range ps.Events {
			f, err := repoClient.GetFile(ctx, m, e.Schema.SchemaRef)
			if err != nil {
				return nil, err
			}
			comp := jsonschema.NewCompiler()
			if err := comp.AddResource(getEventId(m, e.Name), bytes.NewReader(f)); err != nil {
				return nil, err
			}
			sch, err := comp.Compile(getEventId(m, e.Name))
			if err != nil {
				return nil, err
			}

			events[getEventId(m, e.Name)] = &Event{
				event:  e,
				schema: sch,
			}
		}
	}
	return &con, nil
}

func getEventId(merchant string, name string) string {
	return merchant + "/" + name
}

func (rr Controller) getEventDef(merchant string, name string) (*Event, error) {
	q, ok := rr.events[getEventId(merchant, name)]
	if !ok {
		return nil, fmt.Errorf("queue unknown")
	}
	return q, nil
}

func (rr Controller) Routes(logger zerolog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Route("/{merchant_id}/{event_name}", func(r2 chi.Router) {
		r2.Post("/", rr.Post)
		r2.Get("/", rr.Get)
	})
	return r
}

func (rr Controller) Get(w http.ResponseWriter, r *http.Request) {
}

func (c Controller) Post(writer http.ResponseWriter, request *http.Request) {
	logger := zerolog.Ctx(request.Context())
	event := pb.MerchantData{}
	if err := json.NewDecoder(request.Body).Decode(&event); err != nil {
		logger.Error().Msgf("Failed to decode request body: %s", err)
		http.Error(writer, "Request body illegal", http.StatusBadRequest)
		return
	}

	var items []*pb.QueueItem
	merchant := chi.URLParam(request, "merchant_id")
	eventName := chi.URLParam(request, "event_name")
	ef, err := c.getEventDef(merchant, eventName)
	if err != nil {
		logger.Error().Msgf("Data unknown: %s", err)
		http.Error(writer, "Event unknown", http.StatusBadRequest)
	}
	if ef.event.DataType == pb.DataTypes_JSon {

		for _, ev := range event.Data {

			if err = ef.schema.Validate(ev); err != nil {
				logger.Error().Msgf("Failed to validate event: %s", err)
				http.Error(writer, "Failed to validate event", http.StatusBadRequest)
				return
			}

			evBin, err := json.Marshal(ev)
			if err != nil {
				logger.Error().Msgf("Failed to encode event: %s", err)
				http.Error(writer, "Failed to process event", http.StatusBadRequest)
				return
			}
			q := pb.QueueItem{
				Data: string(evBin),
			}
			items = append(items, &q)
		}
	}
	queueRequest := pb.QueueRequest{
		QueueId: &pb.QueueId{
			Name: ef.event.SupplierQueue.GetQueueId(),
		},
		MerchantId: &pb.MerchantId{
			Id: merchant,
		},
		Items: items,
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(c.queueHost, opts...)
	if err != nil {
		logger.Error().Msgf("Failed to connect to queue server: %s", err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	} else {
		client := pb.NewQueueClient(conn)
		defer conn.Close()
		_, err := client.Queue(request.Context(), &queueRequest)
		if err != nil {
			logger.Error().Msgf("Failed to send event to queue server: %s", err)
			http.Error(writer, "", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}

}
