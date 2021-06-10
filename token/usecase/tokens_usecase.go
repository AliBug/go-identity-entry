package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/alibug/go-identity-entry/domain"
	"github.com/alibug/go-identity-utils/config"
	"github.com/alibug/go-identity-utils/status"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// TokensUsecase - 用于操作 token
type TokensUsecase struct {
	tokensRepo  domain.TokensRepository
	tokenConfig config.TokenConfig
}

// NewTokensUsecase will create new an tokenUsecase object representation of domain.TokenUsecase interface
func NewTokensUsecase(repo domain.TokensRepository, tc config.TokenConfig) *TokensUsecase {
	return &TokensUsecase{
		tokensRepo:  repo,
		tokenConfig: tc,
	}
}

// CheckTokensAndLogout - 用于检查 AccessToken 与 RefreshToken 并在存储中删除
func (t *TokensUsecase) CheckTokensAndLogout(ctx context.Context, tokens domain.Tokens) error {
	// 1、CheckTokens
	if tokens.GetAccessToken() != "" {
		// 1.1、检查 AccessToken
		atd, atdExist, err := t.CheckAccessToken(ctx, tokens.GetAccessToken())
		if err != nil {
			return err
		}
		// 1.3、删除 atd
		if atdExist {
			t.deleteTokenID(ctx, atd.GetTokenID())
		}
	}

	if tokens.GetRefreshToken() != "" {
		// 1.2、检查 RefreshToken
		rtd, rtdExist, err := t.CheckRefreshToken(ctx, tokens.GetRefreshToken())
		if err != nil {
			return err
		}
		// 1.4、删除 rtd
		if rtdExist {
			t.deleteTokenID(ctx, rtd.GetTokenID())
		}
	}
	return nil
}

// deleteTokenID - 删除指定 的 Token
func (t *TokensUsecase) deleteTokenID(ctx context.Context, tokenID string) error {
	return t.tokensRepo.DeleteTokenID(ctx, tokenID)
}

// CreateTokens - 同时创建 AccessToken 与 RefreshToken
func (t *TokensUsecase) CreateTokens(ctx context.Context, userID string) (domain.Tokens, error) {
	now := time.Now()
	at, err := t.CreateAccessToken(ctx, userID, now)
	if err != nil {
		return nil, err
	}
	rt, err := t.CreateRefreshToken(ctx, userID, now)
	if err != nil {
		return nil, err
	}
	return &TokensBody{AccessToken: at, RefreshToken: rt}, nil
}

// CreateAccessToken - 创建 AccessToken
func (t *TokensUsecase) CreateAccessToken(ctx context.Context, userID string, now time.Time) (string, error) {
	atUUID := fmt.Sprintf("%s%s%s", uuid.NewString(), "++", userID)
	atParams := NewJwtParams(
		now,
		t.tokenConfig.GetAccessExpirationSeconds(),
		t.tokenConfig.GetAccessTokenSecret(),
		t.tokenConfig.GetIssuer(),
		atUUID,
		userID,
	)
	return t.createToken(ctx, atParams)
}

// CreateRefreshToken - 创建 RefreshToken
func (t *TokensUsecase) CreateRefreshToken(ctx context.Context, userID string, now time.Time) (string, error) {
	rtUUID := fmt.Sprintf("%s%s%s", uuid.NewString(), "++", userID)
	rtParams := NewJwtParams(
		now,
		t.tokenConfig.GetRefreshExpirationSeconds(),
		t.tokenConfig.GetRefreshTokenSecret(),
		t.tokenConfig.GetIssuer(),
		rtUUID,
		userID,
	)
	return t.createToken(ctx, rtParams)
}

// CreateToken - 实现创建 Token
func (t *TokensUsecase) createToken(ctx context.Context, params domain.JwtParams) (string, error) {
	// tokenExpires := params.GetIssueTime().Add(time.Second * params.GetExpiration())
	tokenExpires := params.GetIssueTime().Add(params.GetExpirationSeconds())
	atClaims := jwt.MapClaims{}
	atClaims["aud"] = params.GetAudience()
	atClaims["iss"] = params.GetIssuer()
	atClaims["jti"] = params.GetJwtID()
	atClaims["exp"] = tokenExpires.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atStr, err := at.SignedString(params.GetSecret())
	if err != nil {
		return "", fmt.Errorf("%w : create token error", status.ErrInternalServerError)
	}

	err = t.tokensRepo.CreateTokenID(ctx, NewTokenDetailBody(params.GetJwtID(), params.GetAudience()), params.GetExpirationSeconds())

	if err != nil {
		return "", err
	}
	return atStr, nil
}

// CheckAccessToken - 检查 AccessToken 是否正确
func (t *TokensUsecase) CheckAccessToken(ctx context.Context, tokenStr string) (domain.TokenDetail, bool, error) {
	return t.checkToken(ctx, tokenStr, t.tokenConfig.GetAccessTokenSecret())
}

// CheckRefreshToken - 检查 RefreshToken 是否正确
func (t *TokensUsecase) CheckRefreshToken(ctx context.Context, tokenStr string) (domain.TokenDetail, bool, error) {
	return t.checkToken(ctx, tokenStr, t.tokenConfig.GetRefreshTokenSecret())
}

// checkToken - 返回值 tokenDetail, 是否存在于数据库, 是否出错
func (t *TokensUsecase) checkToken(ctx context.Context, tokenStr string, secret []byte) (domain.TokenDetail, bool, error) {
	td, err := parseJWTToken(tokenStr, secret)
	if err != nil {
		return nil, false, err
	}
	ok, err := t.tokensRepo.CheckTokenID(ctx, td)
	if err != nil {
		return td, false, nil
	}
	if !ok {
		return nil, false, fmt.Errorf("%w : userID not match", status.ErrUnauthorized)
	}
	return td, true, nil
}

func parseJWTToken(tokenStr string, secret []byte) (domain.TokenDetail, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC" SigningMethodHS256
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w : invalid method", status.ErrInternalServerError)
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w : jwt parse error", status.ErrInternalServerError)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		tokenUUID, ok := claims["jti"].(string)
		if !ok {
			return nil, fmt.Errorf("%w : jti not found", status.ErrInternalServerError)
		}
		userID, ok := claims["aud"].(string)
		if err != nil {
			return nil, fmt.Errorf("%w : aud not found", status.ErrInternalServerError)
		}
		return NewTokenDetailBody(tokenUUID, userID), nil
	}
	return nil, err
}
