package cache

import (
	"context"
	"user-service/domain"
)

// FakeCache .
type FakeCache struct{}

// NewFakeCache .
func NewFakeCache() *FakeCache {
	return &FakeCache{}
}

// Set .
func (c *FakeCache) Set(
	ctx context.Context, key string, users []*domain.User) error {
	return nil
}

// Get .
func (c *FakeCache) Get(
	ctx context.Context, key string) ([]*domain.User, bool, error) {
	return nil, false, nil
}
