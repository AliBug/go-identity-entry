package body

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterBody ... 注册信息
type RegisterBody struct {
	Account     string     `json:"account" bson:"account" binding:"required"`
	Password    string     `json:"password"  bson:"-" binding:"required,gte=6"`
	Displayname string     `json:"displayname"  bson:"displayname" binding:"required"`
	CryptPass   []byte     `json:"-" bson:"cryptpass,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// GetAccount - implement domain.RegisterBody
func (r *RegisterBody) GetAccount() string {
	return r.Account
}

// GetPassword - implement domain.RegisterBody
func (r *RegisterBody) GetPassword() string {
	return r.Password
}

// SetCreatedTime - implement domain.RegisterBody
func (r *RegisterBody) SetCreatedTime(t *time.Time) {
	r.CreatedAt = t
}

// SetCryptPass - implement domain.RegisterBody
func (r *RegisterBody) SetCryptPass() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	r.CryptPass = hashedPassword
	return nil
}
