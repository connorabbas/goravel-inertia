package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (r *HomeController) Index(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "Home", map[string]any{
		"message": "Hello from Goravel, Inertia, Vue, and TypeScript.",
	})
}
