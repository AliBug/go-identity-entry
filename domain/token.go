package domain

import (
	"context"
)

// AccessTokenField -
const AccessTokenField = "AccessToken"

// RefreshTokenField -
const RefreshTokenField = "RefreshToken"

// Token contain accessToken and refreshToken
type Token interface {
	GetAccessToken() string
	GetRefreshToken() string
}

// TokenDetail contain tokenID and userID
type TokenDetail interface {
	GetTokenID() string
	GetUserID() string
}

/*
// TokenConfig contain tokens config info
type TokenConfig interface {
	GetAccessTokenSecret() []byte
	GetRefreshTokenSecret() []byte
	GetIssuer() string
	GetAccessExpirationSeconds() time.Duration
	GetRefreshExpirationSeconds() time.Duration
}
*/

// TokenUsecase represent the tokens' usercase contract
type TokenUsecase interface {
	CreateTokenUC(context.Context, string) (Token, error)
	DeleteTokenUC(context.Context, Token) error
	RefreshTokenUC(context.Context, Token) (Token, error)
}

// TokenRepository represent the token' repository cantract
type TokenRepository interface {
	CreateToken(context.Context, string) (Token, error)
	DeleteToken(context.Context, Token) error
	RefreshToken(context.Context, Token) (Token, error)
}
