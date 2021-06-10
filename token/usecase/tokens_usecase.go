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

type tokensUsecase struct {
	tokensRepo  domain.TokensRepository
	tokenConfig config.TokenConfig
}

// NewTokensUsecase will create new an tokenUsecase object representation of domain.TokenUsecase interface
func NewTokensUsecase(repo domain.TokensRepository, tc config.TokenConfig) domain.TokensUseCase {
	return &tokensUsecase{
		tokensRepo:  repo,
		tokenConfig: tc,
	}
}

func (t *tokensUsecase) CheckTokensAndLogout(ctx context.Context, tokens domain.Tokens) error {
	// 1、CheckTokens
	if tokens.GetAccessToken() != "" {
		// 1.1、检查 AccessToken
		atd, atdExist, err := t.CheckAccessToken(ctx, tokens.GetAccessToken())
		if err != nil {
			return err
		}
		// 1.3、删除 atd
		if atdExist {
			t.DeleteTokenID(ctx, atd.GetTokenID())
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
			t.DeleteTokenID(ctx, rtd.GetTokenID())
		}
	}
	return nil
}

// DeleteToken - 删除指定 的 Token
func (t *tokensUsecase) DeleteTokenID(ctx context.Context, tokenID string) error {
	return t.tokensRepo.DeleteTokenID(ctx, tokenID)
}

func (t *tokensUsecase) CreateTokens(ctx context.Context, userID string) (domain.Tokens, error) {
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

func (t *tokensUsecase) CreateAccessToken(ctx context.Context, userID string, now time.Time) (string, error) {
	atUUID := fmt.Sprintf("%s%s%s", uuid.NewString(), "++", userID)
	atParams := NewJwtParams(
		now,
		t.tokenConfig.GetAccessExpirationSeconds(),
		t.tokenConfig.GetAccessTokenSecret(),
		t.tokenConfig.GetIssuer(),
		atUUID,
		userID,
	)
	return t.CreateToken(ctx, atParams)
}

func (t *tokensUsecase) CreateRefreshToken(ctx context.Context, userID string, now time.Time) (string, error) {
	rtUUID := fmt.Sprintf("%s%s%s", uuid.NewString(), "++", userID)
	rtParams := NewJwtParams(
		now,
		t.tokenConfig.GetRefreshExpirationSeconds(),
		t.tokenConfig.GetRefreshTokenSecret(),
		t.tokenConfig.GetIssuer(),
		rtUUID,
		userID,
	)
	return t.CreateToken(ctx, rtParams)
}

// CreateToken - 实现创建 Token
func (t *tokensUsecase) CreateToken(ctx context.Context, params domain.JwtParams) (string, error) {
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

func (t *tokensUsecase) CheckAccessToken(ctx context.Context, tokenStr string) (domain.TokenDetail, bool, error) {
	return t.CheckToken(ctx, tokenStr, t.tokenConfig.GetAccessTokenSecret())
}

func (t *tokensUsecase) CheckRefreshToken(ctx context.Context, tokenStr string) (domain.TokenDetail, bool, error) {
	return t.CheckToken(ctx, tokenStr, t.tokenConfig.GetRefreshTokenSecret())
}

// CheckToken - 返回值 tokenDetail, 是否存在于数据库, 是否出错
func (t *tokensUsecase) CheckToken(ctx context.Context, tokenStr string, secret []byte) (domain.TokenDetail, bool, error) {
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
