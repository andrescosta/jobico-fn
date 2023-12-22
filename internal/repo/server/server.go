package server

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/repo/controller"
)

type Server struct {
	pb.UnimplementedRepoServer
	controller *controller.Controller
}

func New(ctx context.Context, dir string) *Server {
	return &Server{
		controller: controller.New(ctx, dir),
	}
}
func (s *Server) AddFile(ctx context.Context, in *pb.AddFileRequest) (*pb.AddFileReply, error) {
	return s.controller.AddFile(ctx, in)
}
func (s *Server) GetFile(ctx context.Context, in *pb.GetFileRequest) (*pb.GetFileReply, error) {
	return s.controller.GetFile(ctx, in)
}
func (s *Server) GetAllFileNames(ctx context.Context, in *pb.GetAllFileNamesRequest) (*pb.GetAllFileNamesReply, error) {
	return s.controller.GetAllFileNames(ctx, in)
}
func (s *Server) UpdateToFileStr(in *pb.UpdateToFileStrRequest, ctl pb.Repo_UpdateToFileStrServer) error {
	return s.controller.UpdateToFileStr(in, ctl)
}
