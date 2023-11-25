package remote

import (
	"context"
	"io"

	pb "github.com/andrescosta/workflew/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RepoClient struct {
	serverAddr string
}

func NewRepoClient(serverAddr string) *RepoClient {
	return &RepoClient{
		serverAddr: serverAddr,
	}
}

func (c *RepoClient) dial() (*grpc.ClientConn, error) {
	ops := grpc.WithTransportCredentials(insecure.NewCredentials())
	return grpc.Dial(c.serverAddr, ops)

}

func (c *RepoClient) AddFile(ctx context.Context, merchant string, name string, reader io.Reader) error {
	conn, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()
	repo := pb.NewRepoClient(conn)

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = repo.AddFile(ctx, &pb.AddFileRequest{
		MerchantId: merchant,
		Name:       name,
		File:       bytes,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *RepoClient) GetFile(ctx context.Context, merchant string, name string) ([]byte, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	repo := pb.NewRepoClient(conn)
	r, err := repo.GetFile(ctx, &pb.GetFileRequest{
		MerchantId: merchant,
		Name:       name,
	})
	if err != nil {
		return nil, err
	}

	return r.File, nil
}
