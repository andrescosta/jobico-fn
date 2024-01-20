package remote

import (
	"context"
	"io"

	"github.com/andrescosta/goico/pkg/broadcaster"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	rpc "google.golang.org/grpc"
)

type RepoClient struct {
	serverAddr             string
	conn                   *rpc.ClientConn
	client                 pb.RepoClient
	broadcasterRepoUpdates *broadcaster.Broadcaster[*pb.UpdateToFileStrReply]
}

func NewRepoClient(ctx context.Context, d service.GrpcDialer) (*RepoClient, error) {
	addr := env.String("repo.host")
	conn, err := d.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}
	client := pb.NewRepoClient(conn)
	return &RepoClient{
		serverAddr: addr,
		conn:       conn,
		client:     client,
	}, nil
}

func (c *RepoClient) Close() {
	_ = c.conn.Close()
}

func (c *RepoClient) AddFile(ctx context.Context, tenant string, name string, fileType pb.File_FileType, reader io.Reader) error {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	_, err = c.client.AddFile(ctx, &pb.AddFileRequest{
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

func (c *RepoClient) GetFile(ctx context.Context, tenant string, name string) ([]byte, error) {
	r, err := c.client.GetFile(ctx, &pb.GetFileRequest{
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

func (c *RepoClient) GetAllFileNames(ctx context.Context) ([]*pb.TenantFiles, error) {
	reply, err := c.client.GetAllFileNames(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	ret := make([]*pb.TenantFiles, 0)
	ret = append(ret, reply.TenantFiles...)
	return ret, nil
}

func (c *RepoClient) UpdateToFileStr(ctx context.Context, resChan chan<- *pb.UpdateToFileStrReply) error {
	s, err := c.client.UpdateToFileStr(ctx, &pb.UpdateToFileStrRequest{})
	if err != nil {
		return err
	}
	return grpchelper.Recv(ctx, s, resChan)
}

func (c *RepoClient) ListenerForRepoUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToFileStrReply], error) {
	if c.broadcasterRepoUpdates == nil {
		if err := c.startListenRepoUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.broadcasterRepoUpdates.Subscribe()
}

func (c *RepoClient) startListenRepoUpdates(ctx context.Context) error {
	cb := broadcaster.Start[*pb.UpdateToFileStrReply](ctx)
	c.broadcasterRepoUpdates = cb
	s, err := c.client.UpdateToFileStr(ctx, &pb.UpdateToFileStrRequest{})
	if err != nil {
		return err
	}
	go func() {
		_ = grpchelper.Listen(ctx, s, cb)
	}()
	return nil
}
