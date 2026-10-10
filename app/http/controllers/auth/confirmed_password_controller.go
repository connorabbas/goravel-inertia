package auth

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	"goravel/app/http/middleware"
	requests "goravel/app/http/requests/auth"
	"time"
)

type ConfirmedPasswordController struct{}

func (c *ConfirmedPasswordController) Create(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "auth/ConfirmPassword", nil)
}

func (c *ConfirmedPasswordController) Store(ctx http.Context) http.Response {
	var input requests.CurrentPasswordRequest
	if response := Validate(ctx, &input); response != nil {
		return response
	}
	u := middleware.CurrentUser(ctx)
	if !facades.Hash().Check(input.Password, u.Password) {
		return Fail(ctx, "password", "The password is incorrect.")
	}
	ctx.Request().Session().Put("auth.password_confirmed_at", time.Now().Unix())
	return facades.Inertia().Redirect(ctx, middleware.Intended(ctx))
}
