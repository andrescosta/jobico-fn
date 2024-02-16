package client

import (
	"context"
	"io"

	"github.com/andrescosta/goico/pkg/broadcaster"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc/stream"
	pb "github.com/andrescosta/jobico/internal/api/types"
	rpc "google.golang.org/grpc"
)

type Repo struct {
	addr          string
	conn          *rpc.ClientConn
	cli           pb.RepoClient
	bcRepoUpdates *broadcaster.Broadcaster[*pb.UpdateToFileStrReply]
}

func NewRepo(ctx context.Context, d service.GrpcDialer) (*Repo, error) {
	addr := env.String("repo.host")
	conn, err := d.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}
	client := pb.NewRepoClient(conn)
	return &Repo{
		addr: addr,
		conn: conn,
		cli:  client,
	}, nil
}

func (c *Repo) Close() error {
	return c.conn.Close()
}

func (c *Repo) AddFile(ctx context.Context, tenant string, name string, fileType pb.File_FileType, reader io.Reader) error {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	_, err = c.cli.AddFile(ctx, &pb.AddFileRequest{
		TenantFile: &pb.TenantFile{
			Tenant: tenant,
			File: &pb.File{
				Type:    fileType,
				Name:    name,
				Content: bytes,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Repo) File(ctx context.Context, tenant string, name string) ([]byte, error) {
	r, err := c.cli.File(ctx, &pb.FileRequest{
		TenantFile: &pb.TenantFile{
			Tenant: tenant,
			File: &pb.File{
				Name: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return r.File.Content, nil
}

func (c *Repo) AllFilenames(ctx context.Context) ([]*pb.TenantFiles, error) {
	reply, err := c.cli.AllFileNames(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	ret := make([]*pb.TenantFiles, 0)
	ret = append(ret, reply.TenantFiles...)
	return ret, nil
}

func (c *Repo) ListenerForRepoUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToFileStrReply], error) {
	if c.bcRepoUpdates == nil {
		if err := c.startListenRepoUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.bcRepoUpdates.Subscribe()
}

func (c *Repo) startListenRepoUpdates(ctx context.Context) error {
	bc := broadcaster.NewAndStart[*pb.UpdateToFileStrReply](ctx)
	c.bcRepoUpdates = bc
	s, err := c.cli.UpdateToFileStr(ctx, &pb.UpdateToFileStrRequest{})
	if err != nil {
		return err
	}
	go func() {
		_ = stream.Recv(ctx, s, bc)
	}()
	return nil
}
