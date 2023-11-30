package repo

import (
	"context"

	pb "github.com/andrescosta/workflew/api/types"
)

type Server struct {
	pb.UnimplementedRepoServer
	Repo *FileRepo
}

func (s *Server) AddFile(ctx context.Context, r *pb.AddFileRequest) (*pb.AddFileReply, error) {
	err := s.Repo.AddFile(r.TenantId, r.Name, r.File)
	if err != nil {
		return nil, err
	}
	return &pb.AddFileReply{}, nil
}

func (s *Server) GetFile(ctx context.Context, r *pb.GetFileRequest) (*pb.GetFileReply, error) {
	f, err := s.Repo.File(r.TenantId, r.Name)
	if err != nil {
		return nil, err
	}
	return &pb.GetFileReply{
		File: f,
	}, nil
}

func (s *Server) GetAllFileNames(ctx context.Context, r *pb.GetAllFileNamesRequest) (*pb.GetAllFileNamesReply, error) {
	f, err := s.Repo.Files()
	if err != nil {
		return nil, err
	}
	return &pb.GetAllFileNamesReply{
		Files: f,
	}, nil
}
