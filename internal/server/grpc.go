package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rankingpb "gorankd/api/proto/gen"
	"gorankd/internal/ranking"
)

// GRPCServer implements the gRPC ranking service handler.
type GRPCServer struct {
	rankingpb.UnimplementedRankingServiceServer
	ranking ranking.Service
}

// NewGRPCServer creates a new gRPC server handler.
func NewGRPCServer(r ranking.Service) *GRPCServer {
	return &GRPCServer{ranking: r}
}

func (s *GRPCServer) UpdateScore(ctx context.Context, req *rankingpb.UpdateScoreRequest) (*rankingpb.UpdateScoreResponse, error) {
	if err := s.ranking.UpdateScore(ctx, req.PlayerId, req.Score); err != nil {
		return nil, status.Errorf(codes.Internal, "update score: %v", err)
	}
	return &rankingpb.UpdateScoreResponse{}, nil
}

func (s *GRPCServer) GetRank(ctx context.Context, req *rankingpb.GetRankRequest) (*rankingpb.GetRankResponse, error) {
	rank, err := s.ranking.GetRank(ctx, req.PlayerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get rank: %v", err)
	}
	return &rankingpb.GetRankResponse{Rank: rank}, nil
}

func (s *GRPCServer) GetTopN(ctx context.Context, req *rankingpb.GetTopNRequest) (*rankingpb.GetTopNResponse, error) {
	players, err := s.ranking.GetTopN(ctx, int(req.N))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get top n: %v", err)
	}

	pbPlayers := make([]*rankingpb.PlayerRank, len(players))
	for i, p := range players {
		pbPlayers[i] = &rankingpb.PlayerRank{
			PlayerId: p.PlayerID,
			Score:    p.Score,
			Rank:     p.Rank,
		}
	}
	return &rankingpb.GetTopNResponse{Players: pbPlayers}, nil
}

func (s *GRPCServer) GetPlayerScore(ctx context.Context, req *rankingpb.GetPlayerScoreRequest) (*rankingpb.GetPlayerScoreResponse, error) {
	score, err := s.ranking.GetPlayerScore(ctx, req.PlayerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get player score: %v", err)
	}
	return &rankingpb.GetPlayerScoreResponse{Score: score}, nil
}
