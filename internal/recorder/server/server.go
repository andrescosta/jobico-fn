package server

import (
	"context"

	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/recorder/controller"
)

type Server struct {
	pb.UnimplementedRecorderServer
	controller *controller.Recorder
	ctx        context.Context
}

func New(ctx context.Context, fullpath string, o controller.Option) (*Server, error) {
	c, err := controller.New(fullpath, o)
	if err != nil {
		return nil, err
	}
	return &Server{
		controller: c,
		ctx:        ctx,
	}, nil
}

func (s *Server) Close() error {
	return s.controller.Close()
}

func (s *Server) AddJobExecution(_ context.Context, in *pb.AddJobExecutionRequest) (*pb.Void, error) {
	return s.controller.AddJobExecution(s.ctx, in)
}

func (s *Server) GetJobExecutionsStr(in *pb.GetJobExecutionsRequest, srv pb.Recorder_GetJobExecutionsStrServer) error {
	return s.controller.GetJobExecutionsStr(s.ctx, in, srv)
}

func (s *Server) GetJobExecutions(ctx context.Context, in *pb.GetJobExecutionsRequest) (*pb.GetJobExecutionsReply, error) {
	l, err := s.controller.OldRecords(int(*in.Lines))
	if err != nil {
		return nil, err
	}
	return &pb.GetJobExecutionsReply{Result: l}, nil
}
