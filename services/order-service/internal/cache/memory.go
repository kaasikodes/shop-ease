package cache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

// ======================= In-Memory Cache ===========================

type InMemoryCache struct {
	store *cache.Cache
}

func NewInMemoryCache(defaultTTL, cleanupInterval time.Duration) *InMemoryCache {
	return &InMemoryCache{
		store: cache.New(defaultTTL, cleanupInterval),
	}
}

func (c *InMemoryCache) StoreUserInfo(ctx context.Context, user UserInfo) error {
	c.store.SetDefault(userKey(user.Id), user)
	return nil
}

func (c *InMemoryCache) GetUserInfo(ctx context.Context, id int) (*UserInfo, error) {
	data, found := c.store.Get(userKey(id))
	if !found {
		return nil, errors.New("user not found in memory cache")
	}
	info := data.(UserInfo)
	return &info, nil
}

func (c *InMemoryCache) DeleteUserInfo(ctx context.Context, userId int) error {
	c.store.Delete(userKey(userId))
	return nil
}
