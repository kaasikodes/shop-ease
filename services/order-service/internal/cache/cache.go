package cache

import (
	"context"
	"strconv"
)

type UserInfo struct {
	Id    int
	Name  string
	Email string
	Roles []Role
}
type Role struct {
	Id       int
	IsActive bool
	Name     string
}
type CacheRepo interface {
	StoreUserInfo(ctx context.Context, user UserInfo) error
	GetUserInfo(ctx context.Context, id int) (*UserInfo, error)
	DeleteUserInfo(ctx context.Context, userId int) error
}

func userKey(id int) string {
	return "user:" + strconv.Itoa(id)
}
