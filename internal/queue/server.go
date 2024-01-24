package queue

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/queue/controller"
)

type Server struct {
	pb.UnimplementedQueueServer
	controller *controller.Controller
}

func NewServer(ctx context.Context, d service.GrpcDialer, o controller.Option) (*Server, error) {
	c, err := controller.New(ctx, d, o)
	if err != nil {
		return nil, err
	}
	return &Server{
		controller: c,
	}, nil
}

func (s *Server) Queue(ctx context.Context, in *pb.QueueRequest) (*pb.Void, error) {
	return s.controller.Queue(ctx, in)
}

func (s *Server) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	return s.controller.Dequeue(ctx, in)
}
