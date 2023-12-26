package queue

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/queue/controller"
)

type Server struct {
	pb.UnimplementedQueueServer
	controller *controller.Controller
}

func NewServer(ctx context.Context) (*Server, error) {
	c, err := controller.New(ctx)
	if err != nil {
		return nil, err
	}
	return &Server{
		controller: c,
	}, nil
}

func (s *Server) Queue(ctx context.Context, in *pb.QueueRequest) (*pb.QueueReply, error) {
	return s.controller.Queue(ctx, in)
}

func (s *Server) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	return s.controller.Dequeue(ctx, in)
}
