package listener

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Controller struct {
	queueClient *remote.QueueClient
	eventsCache *EventDefCache
}

func New(ctx context.Context, d service.GrpcDialer) (Controller, error) {
	queueClient, err := remote.NewQueueClient(ctx, d)
	var empty Controller
	if err != nil {
		return empty, err
	}
	eventsCache, err := NewCachePopulated(ctx, d)
	if err != nil {
		return empty, err
	}
	return Controller{
		queueClient: queueClient,
		eventsCache: eventsCache,
	}, nil
}

func (c Controller) ConfigureRoutes(_ context.Context, r *mux.Router) error {
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Jobico started."))
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
	ef, err := c.eventsCache.Get(tenant, eventID)
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
	err = c.queueClient.Queue(request.Context(), &queueRequest)
	if err != nil {
		logger.Error().Msgf("Failed to connect to connect to queue server: %s", err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
