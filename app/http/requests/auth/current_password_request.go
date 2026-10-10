package auth

import "github.com/goravel/framework/contracts/http"

type CurrentPasswordRequest struct {
	Password string `form:"password" json:"password"`
}

func (*CurrentPasswordRequest) Authorize(http.Context) error { return nil }

func (*CurrentPasswordRequest) Rules(http.Context) map[string]any {
	return map[string]any{"password": "required|string"}
}
