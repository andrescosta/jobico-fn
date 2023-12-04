package remote

import (
	"context"
	"io"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/grpc"
)

type RepoClient struct {
	serverAddr string
	conn       *grpc.ClientConn
	client     pb.RepoClient
}

func NewRepoClient() (*RepoClient, error) {
	addr := env.GetAsString("repo.host")
	conn, err := service.Dial(addr)
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
	c.conn.Close()
}

func (c *RepoClient) AddFile(ctx context.Context, tenant string, name string, fileType pb.File_FileType, reader io.Reader) error {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = c.client.AddFile(ctx, &pb.AddFileRequest{
		TenantFile: &pb.TenantFile{
			TenantId: tenant,
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
			TenantId: tenant,
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
	reply, err := c.client.GetAllFileNames(ctx, &pb.GetAllFileNamesRequest{})
	if err != nil {
		return nil, err
	}
	ret := make([]*pb.TenantFiles, 0)
	for _, tf := range reply.TenantFiles {
		ret = append(ret, tf)
	}
	return ret, nil
}

func (c *RepoClient) UpdateToFileStr(ctx context.Context, resChan chan<- *pb.UpdateToFileStrReply) error {
	s, err := c.client.UpdateToFileStr(ctx, &pb.UpdateToFileStrRequest{})
	if err != nil {
		return err
	}
	return grpchelper.Recv(ctx, s, resChan)
}
