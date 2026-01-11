package dataplane

import (
	"fmt"
	"log"
	"net"

	pb "github.com/lopster568/phantomDNS/internal/gen/proto/phantomdns/v1"
	"google.golang.org/grpc"
)

// Server wraps the gRPC server.
type Server struct {
	grpcServer *grpc.Server
	port       int
}

// New creates a new gRPC server.
func New(port int, statusSrv pb.DataPlaneStatusServiceServer) *Server {
	s := grpc.NewServer()
	pb.RegisterDataPlaneStatusServiceServer(s, statusSrv)

	return &Server{
		grpcServer: s,
		port:       port,
	}
}

// Start runs the gRPC server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	log.Printf("gRPC server listening on %s", addr)
	return s.grpcServer.Serve(lis)
}

// Stop gracefully stops the server.
func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
