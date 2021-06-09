package redisdb

import (
	"context"
	"time"

	"github.com/alibug/go-identity-entry/domain"
	"github.com/alibug/go-identity-utils/config"
	"github.com/go-redis/redis/v8"
)

type redisTokensRepository struct {
	client *redis.Client
}

// NewRedisTokensRepository will create an object that represent the user.Repository interface
func NewRedisTokensRepository(client *redis.Client, config config.TokenConfig) domain.TokensRepository {
	return &redisTokensRepository{client}
}

// CreateToken - 将指定 tokenID 与 userID 组成键值对 存入数据库
func (r *redisTokensRepository) CreateToken(ctx context.Context, token domain.TokenDetail, expiration time.Duration) error {
	return r.client.Set(ctx, token.GetTokenID(), token.GetUserID(), expiration).Err()
}

// CheckToken - 用给定的 tokenID 对应 userID， 看是否匹配
func (r *redisTokensRepository) CheckToken(ctx context.Context, token domain.TokenDetail) (bool, error) {
	userID, err := r.client.Get(ctx, token.GetTokenID()).Result()
	if err != nil {
		return false, err
	}
	return token.GetUserID() == userID, nil
}

// DeleteToken - 删除指定 的 tokenID
func (r *redisTokensRepository) DeleteToken(ctx context.Context, tokenID string) error {
	return r.client.Del(ctx, tokenID).Err()
}
