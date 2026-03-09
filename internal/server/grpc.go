package server

import (
	"gorankd/internal/ranking"
)

// GRPCServer implements the gRPC ranking service handler.
type GRPCServer struct {
	ranking ranking.Service
}

// NewGRPCServer creates a new gRPC server handler.
func NewGRPCServer(r ranking.Service) *GRPCServer {
	return &GRPCServer{
		ranking: r,
	}
}
