package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/http/middleware"
	"goravel/app/facades"
	"goravel/app/models"
)

type authKey struct{}
type handler struct {
	name string
	run  func(contractshttp.Context)
}

func (h handler) Signature() string { return "app:" + h.name }

func (h handler) Handle(ctx contractshttp.Context) { h.run(ctx) }

func CurrentUser(ctx contractshttp.Context) *models.User {
	u, _ := ctx.Value(authKey{}).(*models.User)
	return u
}

func ResolveUser() contractshttp.Middleware {
	return handler{"resolve_user", func(ctx contractshttp.Context) {
		var user models.User
		if facades.Auth(ctx).Check() && facades.Auth(ctx).User(&user) == nil &&
			ctx.Request().Session().Get("auth_version") == strconv.Itoa(user.AuthVersion) {
			ctx.WithValue(authKey{}, &user)
		} else if cookie := ctx.Request().Cookie("remember_user"); cookie != "" {
			parts := strings.SplitN(cookie, ":", 2)
			if len(parts) == 2 {
				id, err := strconv.ParseUint(parts[0], 10, 64)
				if err == nil && facades.Orm().Query().Find(&user, id) == nil && user.ID != 0 && user.RememberToken != nil && user.RememberExpiresAt != nil && time.Now().Before(*user.RememberExpiresAt) {
					hash := sha256.Sum256([]byte(parts[1]))
					stored, err := hex.DecodeString(*user.RememberToken)
					if err == nil && subtle.ConstantTimeCompare(stored, hash[:]) == 1 {
						if _, err = facades.Auth(ctx).Login(&user); err == nil {
							ctx.Request().Session().Put("auth_version", strconv.Itoa(user.AuthVersion))
							ctx.Response().Cookie(contractshttp.Cookie{Name: "XSRF-TOKEN", Value: ctx.Request().Session().Token(), Path: "/", SameSite: "Lax", Secure: facades.Config().GetBool("session.secure")})
							ctx.WithValue(authKey{}, &user)
						}
					}
				}
			}
		}
		ctx.Request().Next()
	}}
}

// The framework CSRF middleware checks X-CSRF-TOKEN; Inertia sends the token
// from this readable same-site cookie in that header for every mutation.
func CSRF() contractshttp.Middleware {
	verify := middleware.VerifyCsrfToken()
	return handler{"csrf_cookie", func(ctx contractshttp.Context) {
		ctx.Response().Cookie(contractshttp.Cookie{Name: "XSRF-TOKEN", Value: ctx.Request().Session().Token(), Path: "/", SameSite: "Lax", Secure: facades.Config().GetBool("session.secure")})
		verify.Handle(ctx)
	}}
}

func Guest() contractshttp.Middleware {
	return handler{"guest", func(ctx contractshttp.Context) {
		if CurrentUser(ctx) != nil {
			redirect(ctx, "/dashboard")
			return
		}
		ctx.Request().Next()
	}}
}

func Authenticated() contractshttp.Middleware {
	return handler{"authenticated", func(ctx contractshttp.Context) {
		if CurrentUser(ctx) == nil {
			if ctx.Request().Method() == http.MethodGet {
				path := ctx.Request().Path()
				if query := ctx.Request().Origin().URL.RawQuery; query != "" {
					path += "?" + query
				}
				ctx.Request().Session().Put("url.intended", path)
			}
			redirect(ctx, "/login")
			return
		}
		ctx.Request().Next()
	}}
}

func Verified() contractshttp.Middleware {
	return handler{"verified", func(ctx contractshttp.Context) {
		if u := CurrentUser(ctx); u != nil && u.EmailVerifiedAt == nil {
			redirect(ctx, "/verify-email")
			return
		}
		ctx.Request().Next()
	}}
}

func Confirmed() contractshttp.Middleware {
	return handler{"confirmed", func(ctx contractshttp.Context) {
		last, _ := ctx.Request().Session().Get("auth.password_confirmed_at").(int64)
		if time.Now().Unix()-last > 3*60*60 {
			// Only the GET destination is saved: never replay a destructive request.
			ctx.Request().Session().Put("url.intended", "/settings/profile")
			redirect(ctx, "/confirm-password")
			return
		}
		ctx.Request().Next()
	}}
}

func redirect(ctx contractshttp.Context, path string) {
	status := http.StatusFound
	if ctx.Request().Method() != http.MethodGet && ctx.Request().Method() != http.MethodHead {
		status = http.StatusSeeOther
	}
	ctx.Response().Redirect(status, path).Render()
}

func Intended(ctx contractshttp.Context) string {
	value, _ := ctx.Request().Session().Pull("url.intended").(string)
	if value == "" || !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.ContainsAny(value, "\r\n\\") {
		return "/dashboard"
	}
	u, err := url.Parse(value)
	if err != nil || u.Host != "" || u.User != nil {
		return "/dashboard"
	}
	return value
}
