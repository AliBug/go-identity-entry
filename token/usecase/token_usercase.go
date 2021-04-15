package usecase

import (
	"context"
	"time"

	"github.com/alibug/go-identity/domain"
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

func (t *tokenUsecase) CreateTokenUc(c context.Context, userID string) (res domain.Token, err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	res, err = t.tokenRepo.CreateToken(ctx, userID)
	if err != nil {
		return
	}

	return
}

func (t *tokenUsecase) DeleteTokenUc(c context.Context, token domain.Token) (err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	err = t.tokenRepo.DeleteToken(ctx, token)
	if err != nil {
		return
	}

	return
}

func (t *tokenUsecase) RefreshTokenUc(c context.Context, token domain.Token) (res domain.Token, err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	res, err = t.tokenRepo.RefreshToken(ctx, token)
	if err != nil {
		return
	}

	return
}
