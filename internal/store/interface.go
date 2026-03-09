package store

import "context"

// PlayerScore represents a player's score record.
type PlayerScore struct {
	PlayerID string
	Score    float64
}

// Store defines the persistence layer operations.
type Store interface {
	SaveScore(ctx context.Context, playerID string, score float64) error
	GetScore(ctx context.Context, playerID string) (float64, error)
	GetTopScores(ctx context.Context, limit int) ([]PlayerScore, error)
}
