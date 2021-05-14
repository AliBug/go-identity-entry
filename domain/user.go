package domain

import (
	"context"
	"time"
)

// Register ...
type Register interface {
	GetAccount() string
	GetPassword() string
	SetCreatedTime(*time.Time)
	SetCryptPass() error
}

// User ...
type User interface {
	GetUserID() string
	GetDisplayName() string
	GetCryptPass() []byte
	SetUpdatedTime(*time.Time)
}

// UserUsecase ...
type UserUsecase interface {
	RegisterUserUc(ctx context.Context, body Register) error
	GetByIDUc(ctx context.Context, id string) (User, error)
	GetByAccountUc(ctx context.Context, account string) (User, error)
	CheckAccountAndPassUc(ctx context.Context, account string, password string) (User, error)
}

// UserRepository represent the user's repository contract
type UserRepository interface {
	RegisterUser(ctx context.Context, body Register) error
	GetByID(ctx context.Context, id string) (User, error)
	GetByAccount(ctx context.Context, account string) (User, error)
}
