package auth

import "github.com/goravel/framework/contracts/http"

type ResetPasswordRequest struct {
	Email                string `form:"email" json:"email"`
	Token                string `form:"token" json:"token"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

func (*ResetPasswordRequest) Authorize(http.Context) error { return nil }

func (*ResetPasswordRequest) Filters(http.Context) map[string]any {
	return map[string]any{"email": "trim|lower"}
}

func (*ResetPasswordRequest) Rules(http.Context) map[string]any {
	return map[string]any{"email": "required|string|email|max:255", "token": "required|string|size:64", "password": "required|string|min:8|max:255|confirmed", "password_confirmation": "required|string"}
}
