package body

import (
	"time"

	"github.com/alibug/go-identity/domain"
)

// UserBody - Contain Register Body
type UserBody struct {
	ID domain.StrToObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 用户ID
	// RegisterBody
	Account     string     `json:"account" bson:"account" binding:"required"`
	Displayname string     `json:"displayname"  bson:"displayname" binding:"required"`
	CryptPass   []byte     `json:"-" bson:"cryptpass,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"` // 更新时间
}

// GetUserID - implement domain.User
func (u *UserBody) GetUserID() string {
	return string(u.ID)
}

// GetDisplayName - implement domain.User
func (u *UserBody) GetDisplayName() string {
	return u.Displayname
}

// GetCryptPass - implement domain.User
func (u *UserBody) GetCryptPass() []byte {
	return u.CryptPass
}

// SetUpdatedTime - implement domain.User
func (u *UserBody) SetUpdatedTime(t *time.Time) {
	u.UpdatedAt = t
}
