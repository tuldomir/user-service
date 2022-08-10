package cache

import (
	"context"
	"user-service/internal/domain"
)

// FakeEmptyCache .
type FakeEmptyCache struct{}

// NewFakeEmptyCache .
func NewFakeEmptyCache() *FakeEmptyCache {
	return &FakeEmptyCache{}
}

// Set .
func (c *FakeEmptyCache) Set(
	ctx context.Context, key string, users []*domain.User) error {
	return nil
}

// Get .
func (c *FakeEmptyCache) Get(
	ctx context.Context, key string) ([]*domain.User, bool, error) {
	return nil, false, nil
}
