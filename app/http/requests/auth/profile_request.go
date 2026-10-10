package auth

import (
	"errors"
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"goravel/app/http/middleware"
)

type ProfileRequest struct {
	Name  string `form:"name" json:"name"`
	Email string `form:"email" json:"email"`
}

func (*ProfileRequest) Authorize(ctx http.Context) error {
	if middleware.CurrentUser(ctx) == nil {
		return errors.New("unauthenticated")
	}
	return nil
}

func (*ProfileRequest) Filters(http.Context) map[string]any {
	return map[string]any{"name": "trim", "email": "trim|lower"}
}

func (*ProfileRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{"name": "required|string|max:255", "email": fmt.Sprintf("required|string|email|max:255|unique:users,email,id,%d", middleware.CurrentUser(ctx).ID)}
}
