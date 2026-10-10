package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/mail"
	"goravel/app/facades"
	"goravel/app/http/middleware"
	"goravel/app/models"
	"net/url"
	"strconv"
	"time"
)

type EmailVerificationController struct{}

func (c *EmailVerificationController) Create(ctx http.Context) http.Response {
	if middleware.CurrentUser(ctx).EmailVerifiedAt != nil {
		return facades.Inertia().Redirect(ctx, "/dashboard")
	}
	return facades.Inertia().Render(ctx, "auth/VerifyEmail", map[string]any{"status": ctx.Request().Session().Get("status")})
}

func verificationSignature(id uint, email, expires string) string {
	mac := hmac.New(sha256.New, []byte(facades.Config().GetString("app.key")))
	mac.Write([]byte(strconv.FormatUint(uint64(id), 10) + ":" + email + ":" + expires))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *EmailVerificationController) Store(ctx http.Context) http.Response {
	u := middleware.CurrentUser(ctx)
	if u.EmailVerifiedAt != nil {
		return facades.Inertia().Redirect(ctx, "/dashboard")
	}
	if err := SendVerification(u); err != nil {
		facades.Log().Errorf("verification email delivery failed: %v", err)
		return ctx.Response().Status(500).String("Unable to send verification email.")
	}
	ctx.Request().Session().Flash("status", "verification-link-sent")
	return Back(ctx)
}

func SendVerification(u *models.User) error {
	if len(facades.Config().GetString("app.key")) < 32 {
		return errMissingKey
	}
	expires := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	link := facades.Config().GetString("app.url") + "/verify-email/" + strconv.FormatUint(uint64(u.ID), 10) + "/" + verificationSignature(u.ID, u.Email, expires) + "?expires=" + url.QueryEscape(expires)
	return facades.Mail().To([]string{u.Email}).Subject("Verify your email").Content(mail.Content{Text: "Verify your email using this link (valid for one hour): " + link}).Send()
}

func (c *EmailVerificationController) Update(ctx http.Context) http.Response {
	u := middleware.CurrentUser(ctx)
	if len(facades.Config().GetString("app.key")) < 32 {
		return ctx.Response().Status(500).String("Application key is not configured.")
	}
	expires := ctx.Request().Query("expires")
	deadline, err := strconv.ParseInt(expires, 10, 64)
	if err != nil || time.Now().Unix() > deadline || deadline > time.Now().Add(time.Hour).Unix() || ctx.Request().Route("id") != strconv.FormatUint(uint64(u.ID), 10) || !hmac.Equal([]byte(ctx.Request().Route("hash")), []byte(verificationSignature(u.ID, u.Email, expires))) {
		return ctx.Response().Status(403).String("Invalid or expired verification link.")
	}
	if u.EmailVerifiedAt == nil {
		now := time.Now()
		u.EmailVerifiedAt = &now
		if err := facades.Orm().Query().Save(u); err != nil {
			return ctx.Response().Status(500).String("Unable to verify email.")
		}
	}
	return facades.Inertia().Redirect(ctx, "/dashboard")
}
