package controller

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
)

type Controller struct {
	store *Cache[*pb.QueueItem]
}

func New(ctx context.Context, d service.GrpcDialer, o Option) (*Controller, error) {
	s, err := NewQueueCache[*pb.QueueItem](ctx, d, o)
	if err != nil {
		return nil, err
	}
	return &Controller{
		store: s,
	}, nil
}

func (s *Controller) Queue(_ context.Context, in *pb.QueueRequest) (*pb.Void, error) {
	myqueue, err := s.store.GetQueue(in.Tenant, in.Queue)
	if err != nil {
		return nil, err
	}
	for _, i := range in.Items {
		if err := myqueue.Add(i); err != nil {
			return nil, err
		}
	}
	ret := pb.Void{}
	return &ret, nil
}

func (s *Controller) Dequeue(_ context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
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
