package controller

import (
	"context"

	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/repo/provider"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

type Option struct {
	InMemory bool
}

type Controller struct {
	repo        Repository
	bJobPackage *grpchelper.GrpcBroadcaster[*pb.UpdateToFileStrReply, proto.Message]
	ctx         context.Context
}

type Repository interface {
	Add(tenant string, name string, fileType int32, bytes []byte) error
	File(tenant string, name string) ([]byte, error)
	GetMetadataForFile(tenant string, name string) (*provider.Metadata, error)
	Files() ([]*pb.TenantFiles, error)
}

func New(ctx context.Context, dir string, o Option) *Controller {
	r := provider.NewMemRepo()
	if !o.InMemory {
		provider.NewFileRepo(dir)
	}
	return &Controller{
		ctx:         ctx,
		bJobPackage: grpchelper.StartBroadcaster[*pb.UpdateToFileStrReply, proto.Message](ctx),
		repo:        r,
	}
}

func (s *Controller) Close() error {
	return s.bJobPackage.Stop()
}

func (s *Controller) AddFile(ctx context.Context, r *pb.AddFileRequest) (*pb.AddFileReply, error) {
	if err := s.repo.Add(r.TenantFile.Tenant, r.TenantFile.File.Name, int32(r.TenantFile.File.Type), r.TenantFile.File.Content); err != nil {
		return nil, err
	}
	s.bJobPackage.Broadcast(ctx, &pb.TenantFile{Tenant: r.TenantFile.Tenant, File: &pb.File{Name: r.TenantFile.File.Name, Content: r.TenantFile.File.Content}}, pb.UpdateType_New)
	return &pb.AddFileReply{}, nil
}

func (s *Controller) File(_ context.Context, r *pb.FileRequest) (*pb.FileReply, error) {
	f, err := s.repo.File(r.TenantFile.Tenant, r.TenantFile.File.Name)
	if err != nil {
		return nil, err
	}
	m, err := s.repo.GetMetadataForFile(r.TenantFile.Tenant, r.TenantFile.File.Name)
	if err != nil {
		return nil, err
	}
	return &pb.FileReply{
		File: &pb.File{
			Content: f,
			Type:    pb.File_FileType(m.FileType),
		},
	}, nil
}

func (s *Controller) AllFileNames(_ context.Context, _ *pb.Void) (*pb.AllFileNamesReply, error) {
	f, err := s.repo.Files()
	if err != nil {
		return nil, err
	}
	return &pb.AllFileNamesReply{
		TenantFiles: f,
	}, nil
}

func (s *Controller) UpdateToFileStr(_ *pb.UpdateToFileStrRequest, r pb.Repo_UpdateToFileStrServer) error {
	return s.bJobPackage.RcvAndDispatchUpdates(s.ctx, r)
}
