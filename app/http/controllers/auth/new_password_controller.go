package auth

import (
	"crypto/subtle"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	requests "goravel/app/http/requests/auth"
	"goravel/app/models"
	"strconv"
	"time"
)

type NewPasswordController struct{}

func (c *NewPasswordController) Edit(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "auth/ResetPassword", map[string]any{"token": ctx.Request().Route("token"), "email": ctx.Request().Query("email")})
}

func (c *NewPasswordController) Update(ctx http.Context) http.Response {
	var input requests.ResetPasswordRequest
	if response := Validate(ctx, &input); response != nil {
		return response
	}
	email := input.Email
	hash, err := facades.Hash().Make(input.Password)
	if err != nil {
		return ctx.Response().Status(500).String("Unable to reset password.")
	}
	err = facades.Orm().Transaction(func(tx orm.Query) error {
		var reset models.PasswordResetToken
		if err := tx.Where("email", email).LockForUpdate().First(&reset); err != nil {
			return err
		}
		if reset.Email == "" || time.Now().After(reset.ExpiresAt) {
			return errInvalidToken
		}
		var user models.User
		if err := tx.Where("email", email).First(&user); err != nil {
			return err
		}
		if user.ID == 0 {
			return errInvalidToken
		}
		if subtle.ConstantTimeCompare([]byte(reset.TokenHash), []byte(HashToken(input.Token+":"+strconv.FormatUint(uint64(user.ID), 10)))) != 1 {
			return errInvalidToken
		}
		user.Password, user.RememberToken, user.RememberExpiresAt = hash, nil, nil
		user.AuthVersion++
		if err := tx.Save(&user); err != nil {
			return err
		}
		_, err := tx.Delete(&reset)
		return err
	})
	if err == errInvalidToken {
		return Fail(ctx, "email", "This reset link is invalid or expired.")
	}
	if err != nil {
		return ctx.Response().Status(500).String("Unable to reset password.")
	}
	ctx.Request().Session().Flash("status", "Your password has been reset. Please sign in.")
	return facades.Inertia().Redirect(ctx, "/login")
}
