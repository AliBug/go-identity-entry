package usecase

import (
	"context"
	"time"

	"github.com/alibug/go-identity-entry/domain"
)

type tokenUsecase struct {
	tokenRepo      domain.TokenRepository
	contextTimeout time.Duration
}

// NewTokenUsecase will create new an tokenUsecase object representation of domain.TokenUsecase interface
func NewTokenUsecase(repo domain.TokenRepository, timeout time.Duration) domain.TokenUsecase {
	return &tokenUsecase{
		tokenRepo:      repo,
		contextTimeout: timeout,
	}
}

func (t *tokenUsecase) CreateTokenUC(c context.Context, userID string) (res domain.Token, err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	res, err = t.tokenRepo.CreateToken(ctx, userID)
	if err != nil {
		return
	}

	return
}

func (t *tokenUsecase) DeleteTokenUC(c context.Context, token domain.Token) (err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	err = t.tokenRepo.DeleteToken(ctx, token)
	if err != nil {
		return
	}

	return
}

func (t *tokenUsecase) RefreshTokenUC(c context.Context, refreshTokenStr string) (res domain.Token, err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	res, err = t.tokenRepo.RefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return
	}

	return
}
