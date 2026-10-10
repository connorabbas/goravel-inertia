package auth

import "github.com/goravel/framework/contracts/http"

type RegisterRequest struct {
	Name                 string `form:"name" json:"name"`
	Email                string `form:"email" json:"email"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

func (*RegisterRequest) Authorize(http.Context) error { return nil }

func (*RegisterRequest) Filters(http.Context) map[string]any {
	return map[string]any{"name": "trim", "email": "trim|lower"}
}

func (*RegisterRequest) Rules(http.Context) map[string]any {
	return map[string]any{"name": "required|string|max:255", "email": "required|string|email|max:255|unique:users,email", "password": "required|string|min:8|max:255|confirmed", "password_confirmation": "required|string"}
}
