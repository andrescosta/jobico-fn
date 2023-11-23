package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/server"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	server1 "github.com/andrescosta/workflew/ctl/internal/server"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	service.Start(serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	s := grpc.NewServer()
	svr, err := server1.NewQueue(ctx)
	if err != nil {
		return err
	}
	pb.RegisterControlServer(s, svr)
	reflection.Register(s)

	srv, err := server.New(os.Getenv("ctl.port"))
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Info().Msgf("Queue started at:%s", srv.Addr())
	err = srv.ServeGRPC(ctx, s)
	logger.Info().Msg("Queue stopped")
	return err
}

/*	q := &types.QueueDef{
		QueueId: &types.QueueId{
			Name: "aaaa",
		},
	}

	s := Serializer{}
	schema, err := database.Open(context.Background(), ".\\db.db", "queue", &s)
	if err != nil {
		println(err)
	}
	i, err := schema.Add(context.Background(), q)
	if err != nil {
		println(err)
	}
	k, err := schema.Get(context.Background(), i)
	if err != nil {
		println(err)
	}
	println(k.String())
	ks, err := schema.GetAll(context.Background())
	if err != nil {
		println(err)
	}
	for _, kss := range ks {
		fmt.Println(kss)
	}
	schema.Close(context.Background())
*/
//}
