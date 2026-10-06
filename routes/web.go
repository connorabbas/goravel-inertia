package routes

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/http/controllers"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().String(200, "Gonertia: Goravel is running.")
	})

	facades.Route().Static("public", "./public")

	healthController := controllers.NewHealthController()
	// Readiness checks PostgreSQL; the landing page remains independent of it.
	facades.Route().Get("/health/ready", healthController.Index)
}
