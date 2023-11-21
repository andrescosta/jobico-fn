package server

import (
	"context"

	pb "github.com/andrescosta/workflew/api/types"
)

type Server struct {
	pb.UnimplementedQueueServer
}

func (s *Server) GetQueues(ctx context.Context, in *pb.GetQueuesDefRequest) (*pb.GetQueuesDefRequest, error) {
	return nil, nil
}

func (s *Server) AddQueue(ctx context.Context, in *pb.AddQueueDefRequest) (*pb.AddQueueDefReply, error) {
	return nil, nil
}
