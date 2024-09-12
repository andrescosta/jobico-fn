package listener

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Controller struct {
	ctx         context.Context
	queue       *client.Queue
	eventsCache *EventDefCache
}

func NewController(ctx context.Context, d service.GrpcDialer, l service.GrpcListener) (Controller, error) {
	queue, err := client.NewQueue(ctx, d)
	var empty Controller
	if err != nil {
		return empty, err
	}
	eventsCache, err := newCache(ctx, d, l)
	if err != nil {
		return empty, err
	}
	return Controller{
		ctx:         ctx,
		queue:       queue,
		eventsCache: eventsCache,
	}, nil
}

func (c Controller) Close() error {
	var err error
	err = errors.Join(err, c.queue.Close())
	err = errors.Join(err, c.eventsCache.close())
	return err
}

func (c Controller) ConfigureRoutes(_ context.Context, r *mux.Router) error {
	r.HandleFunc("/",
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("Jobico-fn started."))
		}).Methods("GET", "POST")

	s := r.PathPrefix("/events").Subrouter()
	s.Methods("POST").Path("/{tenant_id}/{event_id}").HandlerFunc(c.Post)
	s.Methods("GET").HandlerFunc(c.Get)
	return nil
}

func (c Controller) Get(writer http.ResponseWriter, _ *http.Request) {
	http.Error(writer, "", http.StatusNotFound)
}

func (c Controller) Post(writer http.ResponseWriter, request *http.Request) {
	logger := zerolog.Ctx(request.Context())

	event := pb.MerchantData{}
	if err := json.NewDecoder(request.Body).Decode(&event); err != nil {
		logger.Error().Msgf("Failed to decode request body: %s", err)
		http.Error(writer, "Request body illegal", http.StatusBadRequest)
		return
	}

	tenant := mux.Vars(request)["tenant_id"]
	eventID := mux.Vars(request)["event_id"]
	ef, err := c.eventsCache.Get(c.ctx, tenant, eventID)
	if err != nil {
		logger.Err(err).Msgf("Get: error getting event")
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	if len(event.Data) == 0 {
		http.Error(writer, "Event illegal", http.StatusBadRequest)
		return
	}
	items := make([]*pb.QueueItem, len(event.Data))
	for idx, ev := range event.Data {
		if err = ef.Schema.Validate(ev); err != nil {
			logger.Error().Msgf("Schema.Validate: Failed to validate event: %s", err)
			http.Error(writer, "Event illegal", http.StatusBadRequest)
			return
		}

		evBin, err := json.Marshal(ev)
		if err != nil {
			logger.Error().Msgf("json.Marshal: Failed to encode event: %s", err)
			http.Error(writer, "Failed to process event", http.StatusBadRequest)
			return
		}

		q := pb.QueueItem{
			Data:  evBin,
			Event: eventID,
		}
		items[idx] = &q
	}

	queueRequest := pb.QueueRequest{
		Queue:  ef.EventDef.SupplierQueue,
		Tenant: tenant,
		Items:  items,
	}
	err = c.queue.Queue(request.Context(), &queueRequest)
	if err != nil {
		logger.Error().Msgf("Failed to connect to connect to queue server: %s", err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
