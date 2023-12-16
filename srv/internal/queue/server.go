package queue

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
)

type Server struct {
	pb.UnimplementedQueueServer

	store *Store[*pb.QueueItem]
}

func NewServer(ctx context.Context) (*Server, error) {
	s, err := NewQueueStore[*pb.QueueItem](ctx)

	if err != nil {
		return nil, err
	}

	return &Server{

		store: s,
	}, nil
}

func (s *Server) Queue(_ context.Context, in *pb.QueueRequest) (*pb.QueueReply, error) {
	myqueue, err := s.store.GetQueue(in.Tenant, in.Queue)
	if err != nil {
		return nil, err
	}

	for _, i := range in.Items {
		if err := myqueue.Add(i); err != nil {
			return nil, err
		}
	}
	ret := pb.QueueReply{}
	return &ret, nil
}

func (s *Server) Dequeue(_ context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	myqueue, err := s.store.GetQueue(in.Tenant, in.Queue)

	if err != nil {
		return nil, err
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
