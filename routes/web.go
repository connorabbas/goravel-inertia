package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	sessionmiddleware "github.com/goravel/framework/session/middleware"

	"goravel/app/facades"
	"goravel/app/http/controllers"
	appmiddleware "goravel/app/http/middleware"
	"goravel/app/http/responses"
)

func Web() {
	facades.Route().Static("public", "./public")
	facades.Route().Static("build", "./public/build")

	healthController := controllers.NewHealthController()
	// Readiness checks PostgreSQL; the landing page remains independent of it.
	facades.Route().Get("/health/ready", healthController.Index)
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return responses.Error(ctx, 404)
	})

	facades.Route().Middleware(
		sessionmiddleware.StartSession(),
		appmiddleware.HandleInertiaRequests(),
	).Group(func(router route.Router) {
		router.Get("/", controllers.NewHomeController().Index)
		router.Get("/about", controllers.NewAboutController().Index)
	})
}
