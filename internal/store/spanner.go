package store

import (
	"context"
	"fmt"
)

type spannerStore struct {
	// TODO: add spanner client field
}

// NewSpannerStore creates a new Spanner-backed store.
func NewSpannerStore() Store {
	return &spannerStore{}
}

func (s *spannerStore) SaveScore(_ context.Context, _ string, _ float64) error {
	return fmt.Errorf("not implemented")
}

func (s *spannerStore) GetScore(_ context.Context, _ string) (float64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (s *spannerStore) GetTopScores(_ context.Context, _ int) ([]PlayerScore, error) {
	return nil, fmt.Errorf("not implemented")
}
