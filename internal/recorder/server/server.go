package server

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/recorder/controller"
)

type Server struct {
	pb.UnimplementedRecorderServer
	controller *controller.Recorder
	ctx        context.Context
}

func New(ctx context.Context, fullpath string) (*Server, error) {
	c, err := controller.New(fullpath)
	if err != nil {
		return nil, err
	}
	return &Server{
		controller: c,
		ctx:        ctx,
	}, nil
}

func (s *Server) AddJobExecution(_ context.Context, in *pb.AddJobExecutionRequest) (*pb.AddJobExecutionReply, error) {
	return s.controller.AddJobExecution(s.ctx, in)
}

func (s *Server) GetJobExecutions(in *pb.GetJobExecutionsRequest, srv pb.Recorder_GetJobExecutionsServer) error {
	return s.controller.GetJobExecutions(s.ctx, in, srv)
}
