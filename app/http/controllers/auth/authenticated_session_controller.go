package auth

import (
	"net/http"
	"strconv"

	contractshttp "github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	"goravel/app/http/middleware"
	requests "goravel/app/http/requests/auth"
	"goravel/app/models"
)

type AuthenticatedSessionController struct{}

func NewAuthenticatedSessionController() *AuthenticatedSessionController {
	return &AuthenticatedSessionController{}
}

func (c *AuthenticatedSessionController) Create(ctx contractshttp.Context) contractshttp.Response {
	return facades.Inertia().Render(ctx, "auth/Login", map[string]any{"canRegister": true, "canResetPassword": true, "status": ctx.Request().Session().Get("status")})
}

func (c *AuthenticatedSessionController) Store(ctx contractshttp.Context) contractshttp.Response {
	var input requests.LoginRequest
	if response := Validate(ctx, &input); response != nil {
		return response
	}
	var user models.User
	if err := facades.Orm().Query().Where("email", input.Email).First(&user); err != nil {
		facades.Log().Errorf("login lookup failed: %v", err)
		return ctx.Response().Status(500).String("Unable to sign in.")
	}
	if user.ID == 0 || !facades.Hash().Check(input.Password, user.Password) {
		return Fail(ctx, "email", "These credentials do not match our records.")
	}
	if _, err := facades.Auth(ctx).Login(&user); err != nil {
		return ctx.Response().Status(500).String("Unable to sign in.")
	}
	ctx.Request().Session().Put("auth_version", strconv.Itoa(user.AuthVersion))
	if input.Remember {
		if err := Remember(ctx, &user); err != nil {
			return ctx.Response().Status(500).String("Unable to remember this session.")
		}
	}
	return facades.Inertia().Redirect(ctx, middleware.Intended(ctx))
}

func (c *AuthenticatedSessionController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	var user models.User
	if facades.Auth(ctx).User(&user) == nil && user.RememberToken != nil {
		user.RememberToken, user.RememberExpiresAt = nil, nil
		if err := facades.Orm().Query().Save(&user); err != nil {
			return ctx.Response().Status(500).String("Unable to sign out.")
		}
	}
	if err := facades.Auth(ctx).Logout(); err != nil {
		return ctx.Response().Status(500).String("Unable to sign out.")
	}
	ctx.Request().Session().Flush()
	ctx.Response().Cookie(contractshttp.Cookie{Name: "remember_user", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: "Lax", Secure: facades.Config().GetBool("session.secure")})
	return ctx.Response().Redirect(http.StatusSeeOther, "/")
}
