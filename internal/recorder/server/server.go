package server

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/recorder/controller"
)

type Server struct {
	pb.UnimplementedRecorderServer
	controller *controller.Recorder
}

func New(fullpath string) (*Server, error) {
	c, err := controller.New(fullpath)
	if err != nil {
		return nil, err
	}
	return &Server{
		controller: c,
	}, nil
}

func (s *Server) AddJobExecution(ctx context.Context, in *pb.AddJobExecutionRequest) (*pb.AddJobExecutionReply, error) {
	return s.controller.AddJobExecution(ctx, in)
}
func (s *Server) GetJobExecutions(in *pb.GetJobExecutionsRequest, ctx pb.Recorder_GetJobExecutionsServer) error {
	return s.controller.GetJobExecutions(in, ctx)
}
