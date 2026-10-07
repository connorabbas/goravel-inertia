package responses

import (
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

// Error renders a page for navigation and a toast payload for failed Inertia actions.
// Use this for intentional HTTP errors; validation errors keep their normal flow.
func Error(ctx http.Context, status int) http.Response {
	title := stdhttp.StatusText(status)
	if status == 419 {
		title = "Page Expired"
	}
	config := facades.Config()
	defaults := fmt.Sprintf("errors.defaults.%dxx.", status/100)
	overrides := fmt.Sprintf("errors.statuses.%d.", status)
	detail := config.GetString(overrides+"detail", config.GetString(defaults+"detail"))
	icon := config.GetString(overrides+"icon", config.GetString(defaults+"icon"))

	path := ctx.Request().Path()
	for _, prefix := range []string{"/api", "/health", "/build", "/public"} {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return ctx.Response().Json(status, http.Json{"message": title})
		}
	}
	inertia := ctx.Request().Header("X-Inertia") != ""
	if !inertia && strings.Contains(ctx.Request().Header("Accept"), "application/json") {
		return ctx.Response().Json(status, http.Json{"message": title})
	}
	// Goravel's default recovery logs the cause; keep debug failures out of the UI.
	if status >= 500 && facades.Config().GetBool("app.debug") {
		return ctx.Response().String(status, title)
	}
	method := ctx.Request().Method()
	if inertia && method != "GET" && method != "HEAD" && method != "OPTIONS" {
		return ctx.Response().Json(status, http.Json{
			"status": status, "errorSummary": fmt.Sprintf("%s - %d", title, status),
			"errorDetail": detail, "errorIcon": icon,
		})
	}

	// The adapter has no status argument. Buffer its render to retain the HTTP error
	// code without sending headers early or recreating the Inertia protocol.
	buffer := httptest.NewRecorder()
	wrapped := bufferedContext{requestContext: ctx, writer: buffer}
	facades.Inertia().ShareSession(ctx)
	if err := facades.Inertia().Render(wrapped, "Error", http.Json{
		"status": status, "title": title, "detail": detail, "homepageRoute": "/",
	}).Render(); err != nil {
		facades.Log().Error(err)
		return ctx.Response().String(500, "Internal Server Error")
	}
	for name, values := range buffer.Header() {
		for _, value := range values {
			ctx.Response().Writer().Header().Add(name, value)
		}
	}
	return ctx.Response().Data(status, buffer.Header().Get("Content-Type"), buffer.Body.Bytes())
}

func Recover(ctx http.Context, cause any) {
	facades.Log().WithContext(ctx).Request(ctx.Request()).Error(cause)
	if ctx.Response().Origin().Size() > 0 {
		ctx.Request().Abort()
		return
	}
	if err := Error(ctx, 500).Render(); err != nil {
		facades.Log().Error(err)
	}
	ctx.Request().Abort()
}

type requestContext = http.Context

type bufferedContext struct {
	requestContext
	writer stdhttp.ResponseWriter
}

func (ctx bufferedContext) Response() http.ContextResponse {
	return bufferedResponse{ContextResponse: ctx.requestContext.Response(), writer: ctx.writer}
}

type bufferedResponse struct {
	http.ContextResponse
	writer stdhttp.ResponseWriter
}

func (response bufferedResponse) Writer() stdhttp.ResponseWriter { return response.writer }
