package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/alibug/go-identity/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of domain.ArticleUsecase interface
func NewUserUsecase(a domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       a,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) RegisterUser(c context.Context, body domain.Register) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepo.RegisterUser(ctx, body)
}

func (u *userUsecase) GetByID(c context.Context, id string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (u *userUsecase) GetByUsername(c context.Context, username string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return
	}

	return
}

func (u *userUsecase) CheckUsernameAndPass(c context.Context, username string, password string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w: username or password invalid", domain.ErrBadParamInput)
	}

	// 2、用户存在 则比较密码
	err = bcrypt.CompareHashAndPassword(res.GetCryptPass(), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w: username or password invalid", domain.ErrBadParamInput)
	}
	return res, nil
}
