package domain

import (
	"context"
)

// Tokens contain accessToken and refreshToken
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoginBody ... Account may be username, email, phone number
// Account 可以是 用户名、电子邮件、手机号码
type LoginBody struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required,gte=6"`
}

// TokensUsecase ...
type TokensUsecase interface {
	Login(ctx context.Context, body *LoginBody) (Tokens, error)
}
