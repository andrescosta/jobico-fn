package listener

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrescosta/workflew/api/types"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ListenerController struct {
	Host string
}

func (rr ListenerController) Routes(logger zerolog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Route("/{merchant_id}/{queue_name}", func(r2 chi.Router) {
		r2.Post("/", rr.Post)
		r2.Get("/", rr.Get)
	})
	return r
}

func (rr ListenerController) Get(w http.ResponseWriter, r *http.Request) {
	/*tt := entity.TicketTrans{
		Adventure_id: chi.URLParam(r, "adventure_id"),
		User_id:      chi.URLParam(r, "user_id"),
		Type:         chi.URLParam(r, "type")}
	if ticketTrans, err := rr.service.GetTicketTrans(tt); err != nil {
		rr.logError("Failed to get reservation:", err, r)
		http.Error(w, "Failed to get reservations", http.StatusInternalServerError)
	} else {
		tr := rr.buildTrans(ticketTrans)
		if err = json.NewEncoder(w).Encode(&tr); err != nil {
			rr.logError("Failed to get reservation:", err, r)
			http.Error(w, "Failed to get reservations", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
	*/
}

func (c ListenerController) Post(w http.ResponseWriter, r *http.Request) {

	event := types.MerchantData{}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		c.logError("Failed to decode request body:", err, r)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	var items []*pb.QueueItem

	for _, y := range event.Data {
		bb, err := json.Marshal(y)
		if err != nil {
			c.logError("Failed to decode request body:", err, r)
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}
		q := pb.QueueItem{
			Data: string(bb),
		}
		items = append(items, &q)
	}

	request := pb.QueueRequest{
		QueueId: &pb.QueueId{
			Name: chi.URLParam(r, "queue_name"),
		},
		MerchantId: &pb.MerchantId{
			Id: chi.URLParam(r, "merchant_id"),
		},
		Items: items,
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(c.Host, opts...)

	if err != nil {
		c.logError("Failed to decode request body:", err, r)
	}
	client := pb.NewQueueClient(conn)

	client.Queue(r.Context(), &request)
	defer conn.Close()

	/*tt := entity.TicketTrans{
		Adventure_id: chi.URLParam(r, "adventure_id"),
		User_id:      chi.URLParam(r, "user_id"),
		Type:         chi.URLParam(r, "type"),
	}

	if err := json.NewDecoder(r.Body).Decode(&tt); err != nil {
		rr.logError("Failed to decode request body:", err, r)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	tickets, err := rr.service.GenerateTickets(tt)
	if err != nil {
		rr.logError("Failed to create reservation:", err, r)
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}
	res := rr.buildTrans(tickets)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&res)
	*/
}

func (rr ListenerController) logError(msg string, err error, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	if err != nil {
		msg = fmt.Sprint(msg, err)
	}
	oplog.Error().Msg(msg)
}
