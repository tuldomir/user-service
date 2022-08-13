package cache

import (
	"context"
	"user-service/models"
)

// FakeEmptyCache .
type FakeEmptyCache struct{}

// NewFakeEmptyCache .
func NewFakeEmptyCache() *FakeEmptyCache {
	return &FakeEmptyCache{}
}

// Set .
func (c *FakeEmptyCache) Set(
	ctx context.Context, key string, users []*models.User) error {
	return nil
}

// Get .
func (c *FakeEmptyCache) Get(
	ctx context.Context, key string) ([]*models.User, bool, error) {
	return nil, false, nil
}
