package ranking

import "context"

// PlayerRank represents a player's rank and score.
type PlayerRank struct {
	PlayerID string
	Score    float64
	Rank     int64
}

// Service defines the ranking business logic operations.
type Service interface {
	UpdateScore(ctx context.Context, playerID string, score float64) error
	GetRank(ctx context.Context, playerID string) (int64, error)
	GetTopN(ctx context.Context, n int) ([]PlayerRank, error)
	GetPlayerScore(ctx context.Context, playerID string) (float64, error)
}
