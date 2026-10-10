package auth

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	requests "goravel/app/http/requests/auth"
	"goravel/app/models"
	"strconv"
)

type RegisteredUserController struct{}

func NewRegisteredUserController() *RegisteredUserController { return &RegisteredUserController{} }

func (c *RegisteredUserController) Create(ctx http.Context) http.Response {
	return facades.Inertia().Render(ctx, "auth/Register", nil)
}

func (c *RegisteredUserController) Store(ctx http.Context) http.Response {
	var input requests.RegisterRequest
	if response := Validate(ctx, &input); response != nil {
		return response
	}
	hash, err := facades.Hash().Make(input.Password)
	if err != nil {
		return ctx.Response().Status(500).String("Unable to register.")
	}
	user := models.User{Name: input.Name, Email: input.Email, Password: hash}
	if err := facades.Orm().Query().Create(&user); err != nil {
		facades.Log().Errorf("registration persistence failed: %v", err)
		return ctx.Response().Status(500).String("Unable to register.")
	}
	if _, err := facades.Auth(ctx).Login(&user); err != nil {
		return ctx.Response().Status(500).String("Unable to sign in.")
	}
	ctx.Request().Session().Put("auth_version", strconv.Itoa(user.AuthVersion))
	if err := SendVerification(&user); err != nil {
		facades.Log().Errorf("signup verification email delivery failed: %v", err)
		ctx.Request().Session().Flash("warn_toast", "Unable to send verification email. Please try resending it.")
	}
	return facades.Inertia().Redirect(ctx, "/verify-email")
}
