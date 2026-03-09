package ranking

import (
	"context"
	"fmt"

	"gorankd/internal/cache"
	"gorankd/internal/store"
)

const defaultNamespace = "global"

type rankingService struct {
	cache cache.Cache
	store store.Store
}

// NewService creates a new ranking service.
func NewService(c cache.Cache, s store.Store) Service {
	return &rankingService{
		cache: c,
		store: s,
	}
}

func (s *rankingService) UpdateScore(ctx context.Context, playerID string, score float64) error {
	key := rankKey(defaultNamespace)

	if err := s.cache.ZAdd(ctx, key, playerID, score); err != nil {
		return fmt.Errorf("cache ZAdd: %w", err)
	}

	if err := s.store.SaveScore(ctx, playerID, score); err != nil {
		return fmt.Errorf("store SaveScore: %w", err)
	}

	return nil
}

func (s *rankingService) GetRank(ctx context.Context, playerID string) (int64, error) {
	key := rankKey(defaultNamespace)

	rank, err := s.cache.ZRevRank(ctx, key, playerID)
	if err != nil {
		return 0, fmt.Errorf("cache ZRevRank: %w", err)
	}

	return rank, nil
}

func (s *rankingService) GetTopN(ctx context.Context, n int) ([]PlayerRank, error) {
	key := rankKey(defaultNamespace)

	entries, err := s.cache.ZRevRange(ctx, key, 0, int64(n-1))
	if err != nil {
		return nil, fmt.Errorf("cache ZRevRange: %w", err)
	}

	players := make([]PlayerRank, len(entries))
	for i, e := range entries {
		players[i] = PlayerRank{
			PlayerID: e.Member,
			Score:    e.Score,
			Rank:     int64(i + 1),
		}
	}

	return players, nil
}

func (s *rankingService) GetPlayerScore(ctx context.Context, playerID string) (float64, error) {
	key := rankKey(defaultNamespace)

	score, err := s.cache.ZScore(ctx, key, playerID)
	if err != nil {
		return 0, fmt.Errorf("cache ZScore: %w", err)
	}

	return score, nil
}

func rankKey(namespace string) string {
	return "rank:" + namespace
}
