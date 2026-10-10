package auth

import "github.com/goravel/framework/contracts/http"

type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Remember bool   `form:"remember" json:"remember"`
}

func (*LoginRequest) Authorize(http.Context) error { return nil }

func (*LoginRequest) Filters(http.Context) map[string]any {
	return map[string]any{"email": "trim|lower"}
}

func (*LoginRequest) Rules(http.Context) map[string]any {
	return map[string]any{"email": "required|string|email|max:255", "password": "required|string", "remember": "sometimes|boolean"}
}
