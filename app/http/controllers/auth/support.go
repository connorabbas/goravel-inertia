package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	"goravel/app/models"
)

var errInvalidToken = errors.New("invalid reset token")
var errMissingKey = errors.New("application key is not configured")

func Validate(ctx http.Context, request http.FormRequest) http.Response {
	validationErrors, err := ctx.Request().ValidateRequest(request)
	if err != nil {
		facades.Log().Errorf("request validation failed: %v", err)
		return ctx.Response().Status(500).String("Unable to validate request.")
	}
	if validationErrors != nil {
		facades.Inertia().FlashErrors(ctx, validationErrors)
		return Back(ctx)
	}
	return nil
}

func Fail(ctx http.Context, field, message string) http.Response {
	ctx.Request().Session().Flash("errors", map[string]any{field: message})
	return Back(ctx)
}

func Back(ctx http.Context) http.Response {
	target := "/"
	if ref, err := url.Parse(ctx.Request().Header("Referer")); err == nil && ref.Host == ctx.Request().Host() && strings.HasPrefix(ref.Path, "/") {
		target = ref.RequestURI()
	}
	return facades.Inertia().Redirect(ctx, target)
}

func Random() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func Remember(ctx http.Context, user *models.User) error {
	token, err := Random()
	if err != nil {
		return err
	}
	hash, expires := HashToken(token), time.Now().Add(30*24*time.Hour)
	user.RememberToken, user.RememberExpiresAt = &hash, &expires
	if err := facades.Orm().Query().Save(user); err != nil {
		return err
	}
	ctx.Response().Cookie(http.Cookie{Name: "remember_user", Value: strconv.FormatUint(uint64(user.ID), 10) + ":" + token, Path: "/", SameSite: "Lax", HttpOnly: true, Secure: facades.Config().GetBool("session.secure"), MaxAge: 30 * 24 * 3600})
	return nil
}
