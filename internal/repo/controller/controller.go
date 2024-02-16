package controller

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/broadcaster"
	"github.com/andrescosta/goico/pkg/syncutil"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/repo/provider"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

type Options struct {
	InMemory bool
}

type Controller struct {
	repoProvider Repository
	bJobPackage  *grpchelper.GrpcBroadcaster[*pb.UpdateToFileStrReply, proto.Message]
	ctx          context.Context
	init         *syncutil.OnceDisposable
}

type Repository interface {
	Add(tenant string, name string, fileType int32, bytes []byte) error
	File(tenant string, name string) ([]byte, error)
	GetMetadataForFile(tenant string, name string) (*provider.Metadata, error)
	Files() ([]*pb.TenantFiles, error)
}

func New(ctx context.Context, dir string, o Options) *Controller {
	var repoProvider Repository
	if o.InMemory {
		repoProvider = provider.NewMemRepo()
	} else {
		repoProvider = provider.NewFileRepo(dir)
	}
	return &Controller{
		ctx:          ctx,
		repoProvider: repoProvider,
		init:         syncutil.NewOnceDisposable(),
		bJobPackage:  grpchelper.NewBroadcaster[*pb.UpdateToFileStrReply, proto.Message](ctx),
	}
}

func (s *Controller) Close() error {
	return s.init.Dispose(s.ctx, func(_ context.Context) error {
		err := s.bJobPackage.Stop()
		if errors.Is(err, broadcaster.ErrStopped) {
			return nil
		}
		return err
	})
}

func (s *Controller) AddFile(ctx context.Context, r *pb.AddFileRequest) (*pb.AddFileReply, error) {
	if err := s.repoProvider.Add(r.TenantFile.Tenant, r.TenantFile.File.Name, int32(r.TenantFile.File.Type), r.TenantFile.File.Content); err != nil {
		return nil, err
	}
	err := s.bJobPackage.Broadcast(ctx,
		&pb.TenantFile{
			Tenant: r.TenantFile.Tenant,
			File:   &pb.File{Name: r.TenantFile.File.Name, Content: r.TenantFile.File.Content},
		},
		pb.UpdateType_New)
	if err != nil && !errors.Is(err, broadcaster.ErrStopped) {
		return nil, err
	}
	return &pb.AddFileReply{}, nil
}

func (s *Controller) File(_ context.Context, r *pb.FileRequest) (*pb.FileReply, error) {
	f, err := s.repoProvider.File(r.TenantFile.Tenant, r.TenantFile.File.Name)
	if err != nil {
		return nil, err
	}
	m, err := s.repoProvider.GetMetadataForFile(r.TenantFile.Tenant, r.TenantFile.File.Name)
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
	f, err := s.repoProvider.Files()
	if err != nil {
		return nil, err
	}
	return &pb.AllFileNamesReply{
		TenantFiles: f,
	}, nil
}

func (s *Controller) UpdateToFileStr(_ *pb.UpdateToFileStrRequest, r pb.Repo_UpdateToFileStrServer) error {
	_ = s.init.Do(s.ctx, func(_ context.Context) error {
		s.bJobPackage.Start()
		return nil
	})
	return s.bJobPackage.RcvAndDispatchUpdates(s.ctx, r)
}
