package controllers

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (r *HealthController) Index(ctx http.Context) http.Response {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var result []int
	if err := facades.DB().WithContext(queryContext).Select(&result, "SELECT 1"); err != nil {
		return ctx.Response().Json(503, http.Json{"status": "unavailable"})
	}

	return ctx.Response().Json(200, http.Json{"status": "ok"})
}
