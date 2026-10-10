package auth

import "github.com/goravel/framework/contracts/http"

type PasswordRequest struct {
	CurrentPassword      string `form:"current_password" json:"current_password"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

func (*PasswordRequest) Authorize(http.Context) error { return nil }

func (*PasswordRequest) Rules(http.Context) map[string]any {
	return map[string]any{"current_password": "required|string", "password": "required|string|min:8|max:255|confirmed", "password_confirmation": "required|string"}
}
