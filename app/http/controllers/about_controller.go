package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type AboutController struct{}

func NewAboutController() *AboutController {
	return &AboutController{}
}

func (r *AboutController) Index(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "About", map[string]any{
		"message": "This page is routed by Goravel and visited through Inertia.",
	})
}
