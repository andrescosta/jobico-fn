package queue

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
)

type Server struct {
	pb.UnimplementedQueueServer
	store *QueueStore[*pb.QueueItem]
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

func (s *Server) Queue(ctx context.Context, in *pb.QueueRequest) (*pb.QueueReply, error) {
	myqueue, err := s.store.GetQueue(in.TenantId, in.QueueId)
	if err != nil {
		return nil, err
	}
	for _, i := range in.Items {
		myqueue.Add(i)
	}

	ret := pb.QueueReply{}

	return &ret, nil
}

func (s *Server) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	myqueue, err := s.store.GetQueue(in.TenantId, in.QueueId)
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
