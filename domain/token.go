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
	// CreateTokens - 创建 AccessToken 和 RefreshToken
	CreateTokens(ctx context.Context, userID string) (Tokens, error)

	// CheckTokensAndLogout - 检查 Tokens
	CheckTokensAndLogout(ctx context.Context, tokens Tokens) error

	// CheckAccessToken - 用于检查 AccessToken 合法性
	// CheckAccessToken(ctx context.Context, tokenStr string) (TokenDetail, bool, error)
	// CheckRefreshToken - 用于检查 RefreshToken 合法性
	// CheckRefreshToken(ctx context.Context, tokenStr string) (TokenDetail, bool, error)
	// DeleteToken - 删除指定 的 Token
	// DeleteTokenID(c context.Context, tokenID string) error
}

// TokensRepository - 持久化处理 Tokens
type TokensRepository interface {
	// CreateToken - 创建 指定 TokenDetail
	CreateTokenID(ctx context.Context, token TokenDetail, expiration time.Duration) error
	// CheckAccessToken - 检查 某个 TokenDetail 是否在数据库中持久化保存
	CheckTokenID(ctx context.Context, token TokenDetail) (bool, error)
	// DeleteToken - 删除指定的 Token
	DeleteTokenID(ctx context.Context, tokenID string) error
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
