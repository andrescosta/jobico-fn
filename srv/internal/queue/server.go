package queue

import (
	"context"
	"fmt"

	pb "github.com/andrescosta/jobico/api/types"
)

type Server struct {
	pb.UnimplementedQueueServer
}

func (s *Server) Queue(ctx context.Context, in *pb.QueueRequest) (*pb.QueueReply, error) {
	id := Id{
		QueueId:  in.QueueId,
		TenantId: in.TenantId,
	}
	myqueue, err := GetQueue[*pb.QueueItem](id)
	if err != nil {
		panic(fmt.Sprintf("Error creating directory: %v", err))
	}
	for _, i := range in.Items {
		myqueue.Add(i)
	}

	ret := pb.QueueReply{}

	return &ret, nil
}

func (s *Server) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	id := Id{
		QueueId:  in.QueueId,
		TenantId: in.TenantId,
	}
	myqueue, err := GetQueue[*pb.QueueItem](id)
	if err != nil {
		panic(fmt.Sprintf("Error creating directory: %v", err))
	}
	i, err := myqueue.Remove()
	if err != nil {
		return nil, err
	}
	var iqs []*pb.QueueItem
	if i != nil {
		iqs = append(iqs, i)
	}
	return &pb.DequeueReply{
		Items: iqs,
	}, nil
}
