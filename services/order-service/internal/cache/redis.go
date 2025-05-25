package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ======================= Redis Cache ===========================

type RedisCache struct {
	client    *redis.Client
	keyPrefix string
	ttl       time.Duration
}

func NewRedisCache(addr, password string, db int, prefix string, ttl time.Duration) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db, //redis has 16 logical databases, 0 will be used by default if none is specified, this allows for data segmentation
	})

	return &RedisCache{
		client:    rdb,
		keyPrefix: prefix,
		ttl:       ttl,
	}
}

func (r *RedisCache) key(id int) string {
	return r.keyPrefix + ":" + string(rune(id))
}

func (r *RedisCache) StoreUserInfo(ctx context.Context, user UserInfo) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(user.Id), data, r.ttl).Err()
}

func (r *RedisCache) GetUserInfo(ctx context.Context, id int) (*UserInfo, error) {
	val, err := r.client.Get(ctx, r.key(id)).Result()
	if err == redis.Nil {
		return nil, errors.New("user not found in redis")
	}
	if err != nil {
		return nil, err
	}

	var user UserInfo
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RedisCache) DeleteUserInfo(ctx context.Context, userId int) error {
	return r.client.Del(ctx, userKey(userId)).Err()
}
