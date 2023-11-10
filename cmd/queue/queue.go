package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/internal/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedQueueServer
}

func (s *server) Queue(ctx context.Context, in *pb.QueueRequest) (*pb.QueueReply, error) {
	id := queue.Id{
		Name:     in.QueueId.Name,
		Merchant: in.MerchantId.Id,
	}
	myqueue, err := queue.GetQueue[*pb.QueueItem](id)
	if err != nil {
		panic(fmt.Sprintf("Error creating directory: %v", err))
	}
	for _, i := range in.Items {
		myqueue.Add(i)
	}

	ret := pb.QueueReply{}

	return &ret, nil
}

func (s *server) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	id := queue.Id{
		Name:     in.QueueId.Name,
		Merchant: in.MerchantId.Id,
	}
	myqueue, err := queue.GetQueue[*pb.QueueItem](id)
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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	server := server{}
	pb.RegisterQueueServer(s, &server)
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
