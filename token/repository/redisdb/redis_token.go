package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/alibug/go-identity/domain"
	"github.com/alibug/go-identity/token/repository/body"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type redisTokenRepository struct {
	client      *redis.Client
	tokenConfig domain.TokenConfig
}

// NewRedisTokenRepository will create an object that represent the user.Repository interface
func NewRedisTokenRepository(client *redis.Client, config domain.TokenConfig) domain.TokenRepository {
	return &redisTokenRepository{client, config}
}

func (r *redisTokenRepository) CreateToken(ctx context.Context, userID string) (domain.Token, error) {
	now := time.Now()
	accessTokenExpires := now.Add(time.Second * r.tokenConfig.GetAccessExpirationSeconds())
	accessUUID := uuid.NewString()

	refreshTokenExpires := now.Add(time.Second * r.tokenConfig.GetRefreshExpirationSeconds())
	refreshUUID := fmt.Sprintf("%s%s%s", accessUUID, "++", userID)

	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["aud"] = userID
	atClaims["iss"] = r.tokenConfig.GetIssuer()
	atClaims["jti"] = accessUUID
	atClaims["exp"] = accessTokenExpires.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString(r.tokenConfig.GetAccessTokenSecret())
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["aud"] = userID
	rtClaims["iss"] = r.tokenConfig.GetIssuer()
	rtClaims["jti"] = refreshUUID
	rtClaims["exp"] = refreshTokenExpires.Unix()

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString(r.tokenConfig.GetRefreshTokenSecret())
	if err != nil {
		return nil, err
	}

	// 分别保存了 auth凭证号， 用户ID, 凭证过期时间
	err = r.client.Set(ctx, accessUUID, userID, accessTokenExpires.Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, refreshUUID, userID, refreshTokenExpires.Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	return &body.TokenBody{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type tokenDetail struct {
	tokenUUID string
	userID    string
}

/*
	⚠️ logout 与 delete token 至少有如下 几种情况
	1.1、 cookie 中 accessToken 与 refreshToken 同时存在
	1.2、 cookie 中 只存在 refreshToken
	1.3、 cookie 中 accessToken 与 refreshToken 都不存在， 那就自然退出

	2.1、 app 提交 accessToken + refreshToken 同时存在于 内存
	2.2、 app 提交 aT + rT 中 只有 rT 仍在 redis 中
	2.3、 app 提交 aT + rT 都不存在于 redis 中
*/
func (r *redisTokenRepository) DeleteToken(ctx context.Context, token domain.Token) error {
	// 1、Try to parse access token
	// 	  如果 AccessToken 有效， 则直接按照 accessToken 删除
	if token.GetAccessToken() != "" {
		accessTokenDetail, err := r.CheckAccessToken(ctx, token.GetAccessToken())
		if err != nil {
			// ⚠️ 此处要考虑 GetAccessToken 为空的情况
			return fmt.Errorf("%w:%v", domain.ErrUnauthorized, err)
		}

		// 1.1 不管 token 是否存在 都直接删除 ⚠️ 是否存在 伪造 token 的 可能？
		r.client.Del(ctx, accessTokenDetail.GetTokenID())
		// 1.2 组合出 refreshToken
		refreshUUID := fmt.Sprintf("%s%s%s", accessTokenDetail.GetTokenID(), "++", accessTokenDetail.GetUserID())
		r.client.Del(ctx, refreshUUID)
		return nil
	}

	// ⚠️ 要处理 refresh Token 不为 空 的 情况
	if token.GetRefreshToken() != "" {
		// refreshTokenDetail, err := parseToken(token.GetRefreshToken(), r.tokenConfig.GetRefreshTokenSecret())
		refreshTokenDetail, err := r.checkRefreshToken(ctx, token.GetRefreshToken())
		if err != nil {
			return fmt.Errorf("%w:%v", domain.ErrUnauthorized, err)
		}
		r.client.Del(ctx, refreshTokenDetail.GetTokenID())
		return nil
	}

	return domain.ErrUnauthorized
}

func (r *redisTokenRepository) RefreshToken(ctx context.Context, token domain.Token) (domain.Token, error) {
	if token.GetRefreshToken() != "" {
		// 1、在 AccessToken 失效的情况下， 检查 refresh
		refreshToken, err := r.checkRefreshToken(ctx, token.GetRefreshToken())
		if err != nil {
			return nil, fmt.Errorf("%w:%v", domain.ErrUnauthorized, err)
		}
		// 2、创建 新的 token
		newToken, err := r.CreateToken(ctx, refreshToken.GetUserID())
		if err != nil {
			return nil, fmt.Errorf("%w:%v", domain.ErrInternalServerError, err)
		}
		// 3、删除原有的 refreshToken
		r.client.Del(ctx, refreshToken.GetTokenID())
		return newToken, nil
	}
	return nil, domain.ErrUnauthorized
}

func (r *redisTokenRepository) CheckAccessToken(ctx context.Context, tokenStr string) (domain.TokenDetail, error) {
	return r.checkToken(ctx, tokenStr, r.tokenConfig.GetAccessTokenSecret())
}

func (r *redisTokenRepository) checkRefreshToken(ctx context.Context, tokenStr string) (domain.TokenDetail, error) {
	return r.checkToken(ctx, tokenStr, r.tokenConfig.GetRefreshTokenSecret())
}

func (r *redisTokenRepository) checkToken(ctx context.Context, tokenStr string, secret []byte) (td domain.TokenDetail, err error) {
	td, err = parseToken(tokenStr, secret)
	if err != nil {
		return
	}
	userID, err := r.client.Get(ctx, td.GetTokenID()).Result()
	if err != nil {
		return nil, fmt.Errorf("%w : token expired", domain.ErrInternalServerError)
	}
	if td.GetUserID() != userID {
		return nil, fmt.Errorf("%w : userID not match", domain.ErrInternalServerError)
	}
	return
}

func parseToken(tokenStr string, secret []byte) (domain.TokenDetail, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC" SigningMethodHS256
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w : invalid method", domain.ErrInternalServerError)
		}
		return secret, nil
	})

	if err != nil {
		// ⚠️ 此处可能还有多处 细节值得考虑
		// log.Println("parse err: ", err)
		return nil, fmt.Errorf("%w : jwt parse error", domain.ErrInternalServerError)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		tokenUUID, ok := claims["jti"].(string)
		if !ok {
			return nil, fmt.Errorf("%w : jti not found", domain.ErrInternalServerError)
		}
		userID, ok := claims["aud"].(string)
		if err != nil {
			return nil, fmt.Errorf("%w : aud not found", domain.ErrInternalServerError)
		}
		return body.NewTokenDetailBody(tokenUUID, userID), nil
	}
	return nil, err
}
