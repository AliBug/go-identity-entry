package usecase

import (
	"context"
	"fmt"

	"github.com/alibug/go-identity-entry/domain"
	"github.com/alibug/go-identity-utils/status"
	"github.com/dgrijalva/jwt-go"
)

type tokensUsecase struct {
	tokensRepo domain.TokensRepository
}

// NewTokensUsecase will create new an tokenUsecase object representation of domain.TokenUsecase interface
func NewTokensUsecase(repo domain.TokensRepository) domain.TokensUseCase {
	return &tokensUsecase{
		tokensRepo: repo,
	}
}

// DeleteToken - 删除指定 的 Token
func (t *tokensUsecase) DeleteToken(ctx context.Context, tokenID string) error {
	return t.tokensRepo.DeleteToken(ctx, tokenID)
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

	err = t.tokensRepo.CreateToken(ctx, NewTokenDetailBody(params.GetJwtID(), params.GetAudience()), params.GetExpirationSeconds())

	if err != nil {
		return "", err
	}
	return atStr, nil
}

func (t *tokensUsecase) CheckToken(ctx context.Context, tokenStr string, secret []byte) (domain.TokenDetail, error) {
	td, err := parseJWTToken(tokenStr, secret)
	if err != nil {
		return nil, err
	}
	ok, err := t.tokensRepo.CheckToken(ctx, td)
	if err != nil {
		return nil, fmt.Errorf("%w : token expired", status.ErrUnauthorized)
	}
	if !ok {
		return nil, fmt.Errorf("%w : userID not match", status.ErrUnauthorized)
	}
	return td, nil
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
