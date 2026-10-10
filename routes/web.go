package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/http/limit"
	frameworkmiddleware "github.com/goravel/framework/http/middleware"
	sessionmiddleware "github.com/goravel/framework/session/middleware"

	"goravel/app/facades"
	"goravel/app/http/controllers"
	authcontrollers "goravel/app/http/controllers/auth"
	settingscontrollers "goravel/app/http/controllers/settings"
	appmiddleware "goravel/app/http/middleware"
	"goravel/app/http/responses"
)

func Web() {
	facades.RateLimiter().For("auth", func(ctx http.Context) http.Limit {
		return limit.PerMinute(10).By(ctx.Request().Ip() + ":" + ctx.Request().Path())
	})
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
		appmiddleware.CSRF(),
		appmiddleware.ResolveUser(),
		appmiddleware.HandleInertiaRequests(),
	).Group(func(router route.Router) {
		router.Get("/", controllers.NewHomeController().Index)
		router.Get("/about", controllers.NewAboutController().Index)
		session := authcontrollers.NewAuthenticatedSessionController()
		registered := authcontrollers.NewRegisteredUserController()
		resetLink := &authcontrollers.PasswordResetLinkController{}
		newPassword := &authcontrollers.NewPasswordController{}
		router.Middleware(appmiddleware.Guest()).Group(func(guest route.Router) {
			guest.Get("/login", session.Create)
			guest.Middleware(frameworkmiddleware.Throttle("auth")).Post("/login", session.Store)
			guest.Get("/register", registered.Create)
			guest.Middleware(frameworkmiddleware.Throttle("auth")).Post("/register", registered.Store)
			guest.Get("/forgot-password", resetLink.Create)
			guest.Middleware(frameworkmiddleware.Throttle("auth")).Post("/forgot-password", resetLink.Store)
			guest.Get("/reset-password/{token}", newPassword.Edit)
			guest.Middleware(frameworkmiddleware.Throttle("auth")).Post("/reset-password", newPassword.Update)
		})
		verification := &authcontrollers.EmailVerificationController{}
		confirmed := &authcontrollers.ConfirmedPasswordController{}
		profile := settingscontrollers.NewProfileController()
		password := settingscontrollers.NewPasswordController()
		router.Middleware(appmiddleware.Authenticated()).Group(func(member route.Router) {
			member.Post("/logout", session.Destroy)
			member.Get("/verify-email", verification.Create)
			member.Middleware(frameworkmiddleware.Throttle("auth")).Post("/email/verification-notification", verification.Store)
			member.Get("/verify-email/{id}/{hash}", verification.Update)
			member.Get("/confirm-password", confirmed.Create)
			member.Middleware(frameworkmiddleware.Throttle("auth")).Post("/confirm-password", confirmed.Store)
			member.Get("/settings", func(ctx http.Context) http.Response { return facades.Inertia().Redirect(ctx, "/settings/profile") })
			member.Get("/settings/profile", profile.Edit)
			member.Patch("/settings/profile", profile.Update)
			member.Middleware(appmiddleware.Confirmed()).Delete("/settings/profile", profile.Destroy)
			member.Get("/settings/password", password.Edit)
			member.Put("/settings/password", password.Update)
			member.Get("/settings/appearance", func(ctx http.Context) http.Response { return facades.Inertia().Render(ctx, "settings/Appearance", nil) })
			member.Middleware(appmiddleware.Verified()).Group(func(verified route.Router) {
				verified.Get("/dashboard", func(ctx http.Context) http.Response { return facades.Inertia().Render(ctx, "Dashboard", nil) })
			})
		})
	})
}
