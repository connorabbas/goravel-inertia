package settings

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	auth "goravel/app/http/controllers/auth"
	"goravel/app/http/middleware"
	requests "goravel/app/http/requests/auth"
	"goravel/app/models"
)

type ProfileController struct{}

func NewProfileController() *ProfileController { return &ProfileController{} }

func (c *ProfileController) Edit(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "settings/Profile", map[string]any{"mustVerifyEmail": true, "status": ctx.Request().Session().Get("status")})
}

func (c *ProfileController) Update(ctx http.Context) http.Response {
	var input requests.ProfileRequest
	if response := auth.Validate(ctx, &input); response != nil {
		return response
	}
	u := middleware.CurrentUser(ctx)
	u.Name = input.Name
	if u.Email != input.Email {
		if _, err := facades.Orm().Query().Where("email", u.Email).Delete(&models.PasswordResetToken{}); err != nil {
			return ctx.Response().Status(500).String("Unable to save profile.")
		}
		u.Email, u.EmailVerifiedAt = input.Email, nil
	}
	if err := facades.Orm().Query().Save(u); err != nil {
		facades.Log().Errorf("profile persistence failed: %v", err)
		return ctx.Response().Status(500).String("Unable to save profile.")
	}
	ctx.Request().Session().Flash("success_toast", "Profile updated.")
	return auth.Back(ctx)
}

func (c *ProfileController) Destroy(ctx http.Context) http.Response {
	var input requests.CurrentPasswordRequest
	if response := auth.Validate(ctx, &input); response != nil {
		return response
	}
	u := middleware.CurrentUser(ctx)
	if !facades.Hash().Check(input.Password, u.Password) {
		return auth.Fail(ctx, "password", "The password is incorrect.")
	}
	if _, err := facades.Orm().Query().Where("email", u.Email).Delete(&models.PasswordResetToken{}); err != nil {
		return ctx.Response().Status(500).String("Unable to delete account.")
	}
	if _, err := facades.Orm().Query().Delete(u); err != nil {
		return ctx.Response().Status(500).String("Unable to delete account.")
	}
	if err := facades.Auth(ctx).Logout(); err != nil {
		return ctx.Response().Status(500).String("Unable to sign out.")
	}
	ctx.Request().Session().Flush()
	ctx.Response().Cookie(http.Cookie{Name: "remember_user", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: "Lax", Secure: facades.Config().GetBool("session.secure")})
	return facades.Inertia().Redirect(ctx, "/")
}
