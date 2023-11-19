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

type Controller struct {
	QueueHost string
}

func (rr Controller) Routes(logger zerolog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Route("/{merchant_id}/{queue_name}", func(r2 chi.Router) {
		r2.Post("/", rr.Post)
		r2.Get("/", rr.Get)
	})
	return r
}

func (rr Controller) Get(w http.ResponseWriter, r *http.Request) {
}

func (c Controller) Post(w http.ResponseWriter, r *http.Request) {

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
	conn, err := grpc.Dial(c.QueueHost, opts...)

	if err != nil {
		c.logError("Failed to decode request body:", err, r)
	}
	client := pb.NewQueueClient(conn)

	client.Queue(r.Context(), &request)
	defer conn.Close()

}

func (rr Controller) logError(msg string, err error, r *http.Request) {
	oplog := zerolog.Ctx(r.Context())
	if err != nil {
		msg = fmt.Sprint(msg, err)
	}
	oplog.Error().Msg(msg)
}
