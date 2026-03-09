package cache

import "context"

// RankEntry represents a member and score pair from a sorted set.
type RankEntry struct {
	Member string
	Score  float64
}

// Cache defines the caching layer operations for ranking.
type Cache interface {
	ZAdd(ctx context.Context, key string, member string, score float64) error
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]RankEntry, error)
	ZRevRank(ctx context.Context, key string, member string) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
}
