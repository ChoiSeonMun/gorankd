package cache

import (
	"context"
	"fmt"
)

type redisCache struct {
	// TODO: add redis client field
}

// NewRedisCache creates a new Redis-backed cache.
func NewRedisCache() Cache {
	return &redisCache{}
}

func (r *redisCache) ZAdd(_ context.Context, _ string, _ string, _ float64) error {
	return fmt.Errorf("not implemented")
}

func (r *redisCache) ZRevRange(_ context.Context, _ string, _, _ int64) ([]RankEntry, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *redisCache) ZRevRank(_ context.Context, _ string, _ string) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (r *redisCache) ZScore(_ context.Context, _ string, _ string) (float64, error) {
	return 0, fmt.Errorf("not implemented")
}
