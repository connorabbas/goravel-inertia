package auth

import "github.com/goravel/framework/contracts/http"

type EmailRequest struct {
	Email string `form:"email" json:"email"`
}

func (*EmailRequest) Authorize(http.Context) error { return nil }

func (*EmailRequest) Filters(http.Context) map[string]any {
	return map[string]any{"email": "trim|lower"}
}

func (*EmailRequest) Rules(http.Context) map[string]any {
	return map[string]any{"email": "required|string|email|max:255"}
}
