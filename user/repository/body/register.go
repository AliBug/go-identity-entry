package body

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterBody ... 注册信息
type RegisterBody struct {
	Username    string     `json:"username" bson:"username" binding:"required"`
	Password    string     `json:"password"  bson:"-" binding:"required,gte=6"`
	Displayname string     `json:"displayname"  bson:"displayname" binding:"required"`
	CryptPass   []byte     `json:"-" bson:"cryptpass,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// GetUsername - implement domain.RegisterBody
func (r *RegisterBody) GetUsername() string {
	return r.Username
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
