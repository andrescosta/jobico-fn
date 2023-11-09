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
	myqueue queue.Queue[*pb.QueueItem]
}

func (s *server) Queue(ctx context.Context, in *pb.QueueRequest) (*pb.QueueReply, error) {
	ret := pb.QueueReply{
		Result: &pb.Result{
			Code: "0",
		},
	}
	for _, i := range in.Items {
		s.myqueue.Add(i)
	}

	return &ret, nil
}

func (s *server) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	i, _ := s.myqueue.Remove()
	var iqs []*pb.QueueItem
	if i == nil {
		return &pb.DequeueReply{
			Result: &pb.Result{Message: "Empty"},
		}, nil

	}
	iqs = append(iqs, i)
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
	server := server{
		myqueue: queue.GetDefault[*pb.QueueItem](),
	}
	pb.RegisterQueueServer(s, &server)
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
