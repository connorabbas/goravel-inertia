package auth

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/mail"
	"goravel/app/facades"
	requests "goravel/app/http/requests/auth"
	"goravel/app/models"
	"net/url"
	"strconv"
	"time"
)

type PasswordResetLinkController struct{}

func (c *PasswordResetLinkController) Create(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "auth/ForgotPassword", map[string]any{"status": ctx.Request().Session().Get("status")})
}

func (c *PasswordResetLinkController) Store(ctx http.Context) http.Response {
	var input requests.EmailRequest
	if response := Validate(ctx, &input); response != nil {
		return response
	}
	email := input.Email
	var user models.User
	if err := facades.Orm().Query().Where("email", email).First(&user); err != nil {
		return ctx.Response().Status(500).String("Unable to request reset.")
	}
	if user.ID != 0 {
		token, err := Random()
		if err != nil {
			return ctx.Response().Status(500).String("Unable to request reset.")
		}
		_, err = facades.Orm().Query().Where("email", email).Delete(&models.PasswordResetToken{})
		if err != nil {
			return ctx.Response().Status(500).String("Unable to request reset.")
		}
		if err := facades.Orm().Query().Create(&models.PasswordResetToken{Email: email, TokenHash: HashToken(token + ":" + strconv.FormatUint(uint64(user.ID), 10)), ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
			return ctx.Response().Status(500).String("Unable to request reset.")
		}
		link := facades.Config().GetString("app.url") + "/reset-password/" + token + "?email=" + url.QueryEscape(email)
		if err := facades.Mail().To([]string{email}).Subject("Reset your password").Content(mail.Content{Text: "Reset your password using this link (valid for one hour): " + link}).Send(); err != nil {
			facades.Log().Errorf("password reset email delivery failed: %v", err)
			return ctx.Response().Status(500).String("Unable to send reset email.")
		}
	}
	ctx.Request().Session().Flash("status", "If that account exists, a password reset link has been sent.")
	return Back(ctx)
}
