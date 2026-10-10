package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

type User struct {
	orm.Model
	Name              string
	Email             string `gorm:"uniqueIndex"`
	Password          string
	EmailVerifiedAt   *time.Time
	RememberToken     *string
	RememberExpiresAt *time.Time
	AuthVersion       int
}

func (u User) Public() map[string]any {
	return map[string]any{"id": u.ID, "name": u.Name, "email": u.Email,
		"emailVerifiedAt": u.EmailVerifiedAt, "createdAt": u.CreatedAt, "updatedAt": u.UpdatedAt}
}
