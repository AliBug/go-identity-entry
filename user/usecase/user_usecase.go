package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/alibug/go-identity-entry/domain"
	"github.com/alibug/go-identity-utils/status"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of domain.ArticleUsecase interface
func NewUserUsecase(repo domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       repo,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) RegisterUserUC(c context.Context, body domain.Register) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepo.RegisterUser(ctx, body)
}

func (u *userUsecase) GetByIDUC(c context.Context, id string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (u *userUsecase) GetByAccountUC(c context.Context, username string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.GetByAccount(ctx, username)
	if err != nil {
		return
	}

	return
}

func (u *userUsecase) CheckAccountAndPassUC(c context.Context, username string, password string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.GetByAccount(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w: username or password invalid", status.ErrBadParamInput)
	}

	// 2、用户存在 则比较密码
	err = bcrypt.CompareHashAndPassword(res.GetCryptPass(), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w: username or password invalid", status.ErrBadParamInput)
	}
	return res, nil
}
