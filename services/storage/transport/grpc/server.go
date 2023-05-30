package grpc

import (
	"fmt"
	"net"

	pb "github.com/rvinnie/lightstream/services/storage/pb"
	"google.golang.org/grpc"
)

type Server struct {
	imageStorageHandler pb.ImageStorageServer
	srv                 *grpc.Server
}

func NewServer(imageStorageHandler pb.ImageStorageServer) *Server {
	return &Server{
		srv:                 grpc.NewServer(),
		imageStorageHandler: imageStorageHandler,
	}
}

func (s *Server) ListenAndServe(port string) error {
	addr := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	pb.RegisterImageStorageServer(s.srv, s.imageStorageHandler)

	if err = s.srv.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
