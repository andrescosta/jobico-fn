package listener

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Controller struct {
	queueClient *remote.QueueClient
	eventsCache *EventDefCache
}

func ConfigureRoutes(ctx context.Context, r *mux.Router) error {
	queueClient, err := remote.NewQueueClient(ctx)
	if err != nil {
		return err
	}
	eventsCache, err := NewEventDefCache(ctx)
	if err != nil {
		return err
	}
	c := Controller{
		queueClient: queueClient,
		eventsCache: eventsCache,
	}
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

	var items []*pb.QueueItem
	tenant := mux.Vars(request)["tenant_id"]
	eventID := mux.Vars(request)["event_id"]
	ef, err := c.eventsCache.Get(tenant, eventID)
	if err != nil {
		logger.Err(err).Msgf("Get: error getting event")
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}

	if ef.EventDef.DataType == pb.DataType_Json {
		for _, ev := range event.Data {
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
			items = append(items, &q)
		}
	}
	queueRequest := pb.QueueRequest{
		Queue:  ef.EventDef.SupplierQueue,
		Tenant: tenant,
		Items:  items,
	}
	err = c.queueClient.Queue(request.Context(), &queueRequest)
	if err != nil {
		logger.Error().Msgf("Failed to connect to connect to queue server: %s", err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
