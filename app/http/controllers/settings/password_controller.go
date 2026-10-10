package settings

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	auth "goravel/app/http/controllers/auth"
	"goravel/app/http/middleware"
	requests "goravel/app/http/requests/auth"
	"strconv"
)

type PasswordController struct{}

func NewPasswordController() *PasswordController { return &PasswordController{} }

func (c *PasswordController) Edit(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "settings/Password", nil)
}

func (c *PasswordController) Update(ctx http.Context) http.Response {
	var input requests.PasswordRequest
	if response := auth.Validate(ctx, &input); response != nil {
		return response
	}
	u := middleware.CurrentUser(ctx)
	if !facades.Hash().Check(input.CurrentPassword, u.Password) {
		return auth.Fail(ctx, "current_password", "The current password is incorrect.")
	}
	hash, err := facades.Hash().Make(input.Password)
	if err != nil {
		return ctx.Response().Status(500).String("Unable to update password.")
	}
	u.Password, u.RememberToken, u.RememberExpiresAt = hash, nil, nil
	u.AuthVersion++
	if err := facades.Orm().Query().Save(u); err != nil {
		return ctx.Response().Status(500).String("Unable to update password.")
	}
	ctx.Request().Session().Put("auth_version", strconv.Itoa(u.AuthVersion))
	ctx.Response().Cookie(http.Cookie{Name: "remember_user", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: "Lax", Secure: facades.Config().GetBool("session.secure")})
	ctx.Request().Session().Flash("success_toast", "Password updated.")
	return auth.Back(ctx)
}
