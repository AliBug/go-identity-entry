package domain

import (
	"context"
	"time"
)

// TokenDetail contain tokenID and userID
type TokenDetail interface {
	GetTokenID() string
	GetUserID() string
}

// Tokens - 包含 AccessToken 和 RefreshToken
type Tokens interface {
	GetAccessToken() string
	GetRefreshToken() string
}

// TokensUseCase - 处理 Tokens
type TokensUseCase interface {
	// CreateToken - 创建 token
	CreateToken(ctx context.Context, params JwtParams) (string, error)

	// CheckToken - 用于检查 token 合法性
	CheckToken(c context.Context, tokenStr string, secret []byte) (TokenDetail, error)

	// DeleteToken - 删除指定 的 Token
	DeleteToken(c context.Context, tokenID string) error
}

// TokensRepository - 持久化处理 Tokens
type TokensRepository interface {
	// CreateToken - 创建 指定 TokenDetail
	CreateToken(ctx context.Context, token TokenDetail, expiration time.Duration) error
	// CheckAccessToken - 检查 某个 TokenDetail 是否在数据库中持久化保存
	CheckToken(ctx context.Context, token TokenDetail) (bool, error)
	// DeleteToken - 删除指定的 Token
	DeleteToken(ctx context.Context, tokenID string) error
}

// JwtParams - 创建 JWT 要用的参数
type JwtParams interface {
	GetExpirationSeconds() time.Duration
	GetIssuer() string
	GetJwtID() string
	GetAudience() string
	GetSecret() []byte
	GetIssueTime() time.Time
}
