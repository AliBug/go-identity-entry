package domain

import (
	"context"
	"time"
)

// Tokens contain accessToken and refreshToken
/*
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
*/

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

// TokenConfig contain tokens config info
type TokenConfig interface {
	GetAccessTokenSecret() []byte
	GetRefreshTokenSecret() []byte
	GetIssuer() string
	GetAccessExpirationSeconds() time.Duration
	GetRefreshExpirationSeconds() time.Duration
}

// TokenUsecase represent the tokens' usercase contract
type TokenUsecase interface {
	CreateTokenUc(context.Context, string) (Token, error)
	DeleteTokenUc(context.Context, Token) error
	RefreshTokenUc(context.Context, Token) (Token, error)
}

// TokenRepository represent the token' repository cantract
type TokenRepository interface {
	CreateToken(context.Context, string) (Token, error)
	DeleteToken(context.Context, Token) error
	RefreshToken(context.Context, Token) (Token, error)
}

// LoginBody ... Account may be username, email, phone number
// Account 可以是 用户名、电子邮件、手机号码
/*
type LoginBody struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required,gte=6"`
}
*/

// TokensUsecase ...
/*
type TokensUsecase interface {
	Login(ctx context.Context, body *LoginBody) (Tokens, error)
}
*/
