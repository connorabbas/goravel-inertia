package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/foundation"
	sessionmiddleware "github.com/goravel/framework/session/middleware"
	frameworkmock "github.com/goravel/framework/testing/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goravel/app/facades"
	"goravel/app/http/responses"
	"goravel/bootstrap"
)

func TestMain(m *testing.M) {
	if name := os.Getenv("GONERTIA_TEST_DATABASE"); name != "" {
		config := facades.Config()
		if !strings.HasPrefix(name, "gonertia_test_") || name == config.GetString("database.connections.postgres.database") {
			fmt.Fprintln(os.Stderr, "GONERTIA_TEST_DATABASE must be a separate gonertia_test_* database")
			os.Exit(1)
		}
		// Providers cache connections during boot, so select the test database first.
		config.Add("database.connections.postgres.database", name)
	}
	bootstrap.Boot()
	os.Exit(m.Run())
}

func TestHTTP(t *testing.T) {
	router := facades.Route()
	app := foundation.App
	defer func() { foundation.App = app }()

	for _, test := range []struct {
		name string
		err  error
		code int
		body string
	}{
		{"ready", nil, 200, `{"status":"ok"}`},
		{"database unavailable", errors.New("private database credentials"), 503, `{"status":"unavailable"}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			database := frameworkmock.Factory().DB()
			database.EXPECT().WithContext(mock.Anything).Return(database).Once()
			database.EXPECT().Select(mock.Anything, "SELECT 1").Return(test.err).Once()
			response := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/health/ready", nil)
			request.Header.Set("X-Inertia", "true")
			request.Header.Set("X-Inertia-Version", "stale")
			router.ServeHTTP(response, request)
			require.Equal(t, test.code, response.Code)
			require.JSONEq(t, test.body, response.Body.String())
			require.Empty(t, response.Header().Get("X-Inertia"))
			require.Empty(t, response.Header().Get("Set-Cookie"))
			database.AssertExpectations(t)
		})
	}
}

func TestInertiaPages(t *testing.T) {
	for _, test := range []struct {
		path, component, message string
	}{
		{"/", "Home", "Hello from Goravel, Inertia, Vue, and TypeScript."},
		{"/about", "About", "This page is routed by Goravel and visited through Inertia."},
	} {
		t.Run(test.component, func(t *testing.T) {
			response := httptest.NewRecorder()
			facades.Route().ServeHTTP(response, httptest.NewRequest("GET", test.path, nil))
			require.Equal(t, 200, response.Code)
			require.Contains(t, response.Header().Get("Content-Type"), "text/html")
			require.Contains(t, response.Body.String(), `<div id="app" class="isolate"></div>`)
			_, embedded, found := strings.Cut(response.Body.String(), `<script data-page="app" type="application/json">`)
			require.True(t, found)
			embedded, _, found = strings.Cut(embedded, "</script>")
			require.True(t, found)
			var initial struct {
				Component string
				URL       string
				Version   string
				Props     map[string]any
			}
			require.NoError(t, json.Unmarshal([]byte(embedded), &initial))
			require.Equal(t, test.component, initial.Component)
			require.Equal(t, test.path, initial.URL)
			require.Equal(t, test.message, initial.Props["message"])
			require.Equal(t, facades.Config().GetString("app.name"), initial.Props["appName"])

			request := httptest.NewRequest("GET", test.path, nil)
			request.Header.Set("X-Inertia", "true")
			request.Header.Set("X-Inertia-Version", initial.Version)
			response = httptest.NewRecorder()
			facades.Route().ServeHTTP(response, request)
			require.Equal(t, 200, response.Code)
			require.Equal(t, "true", response.Header().Get("X-Inertia"))
			require.Contains(t, response.Header().Get("Content-Type"), "application/json")
			var visit map[string]any
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &visit))
			require.Equal(t, test.component, visit["component"])
			require.Equal(t, test.path, visit["url"])
			require.Equal(t, test.message, visit["props"].(map[string]any)["message"])

			request.Header.Set("X-Inertia-Version", initial.Version+"-stale")
			response = httptest.NewRecorder()
			facades.Route().ServeHTTP(response, request)
			require.Equal(t, 409, response.Code)
			require.NotEmpty(t, response.Header().Get("X-Inertia-Location"))
		})
	}
}

func TestErrorResponses(t *testing.T) {
	config := facades.Config()
	debug := config.GetBool("app.debug")
	config.Add("app.debug", false)
	t.Cleanup(func() { config.Add("app.debug", debug) })
	router := facades.Route()
	for _, status := range []int{401, 403, 404, 405, 419, 429, 500, 502, 503} {
		path := fmt.Sprintf("/_test/error/%d", status)
		handler := func(ctx http.Context) http.Response { return responses.Error(ctx, status) }
		router.Get(path, handler)
		router.Post(path, handler)
	}
	panicHandler := func(ctx http.Context) http.Response { panic("private database password") }
	router.Get("/_test/panic", panicHandler)
	router.Post("/_test/panic", panicHandler)
	router.Get("/api/_test/panic", panicHandler)

	for _, test := range []struct {
		path string
		code int
	}{
		{"/missing-page", 404}, {"/_test/error/403", 403},
		{"/_test/error/500", 500}, {"/_test/error/503", 503}, {"/_test/panic", 500},
	} {
		for _, inertia := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/inertia=%t", test.path, inertia), func(t *testing.T) {
				request := httptest.NewRequest("GET", test.path, nil)
				if inertia {
					request.Header.Set("X-Inertia", "true")
					request.Header.Set("X-Inertia-Version", facades.Inertia().Version())
				}
				response := httptest.NewRecorder()
				router.ServeHTTP(response, request)
				require.Equal(t, test.code, response.Code)
				require.NotContains(t, response.Body.String(), "private database password")
				body := response.Body.String()
				if !inertia {
					_, body, _ = strings.Cut(body, `<script data-page="app" type="application/json">`)
					body, _, _ = strings.Cut(body, "</script>")
				}
				var page map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &page))
				require.Equal(t, "Error", page["component"])
				props := page["props"].(map[string]any)
				require.Equal(t, float64(test.code), props["status"])
				require.Equal(t, "/", props["homepageRoute"])
				require.Equal(t, config.GetString("app.name"), props["appName"])
			})
		}
	}

	t.Run("failed mutation returns a toast payload", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/_test/error/403", nil)
		request.Header.Set("X-Inertia", "true")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		require.Equal(t, 403, response.Code)
		require.Empty(t, response.Header().Get("X-Inertia"))
		var payload map[string]any
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &payload))
		require.Equal(t, float64(403), payload["status"])
		require.NotEmpty(t, payload["errorSummary"])
		require.Equal(t, "Sorry, you are unauthorized to access this resource/action.", payload["errorDetail"])
		require.Equal(t, "i-lucide-shield-alert", payload["errorIcon"])
	})
	t.Run("configured status messages and family defaults", func(t *testing.T) {
		for _, test := range []struct {
			status int
			detail string
			icon   string
		}{
			{401, "Please sign in to continue.", "i-lucide-log-in"},
			{404, "Sorry, the resource you are looking for could not be found.", "i-lucide-search-x"},
			{405, "Sorry, your request could not be completed.", "i-lucide-circle-alert"},
			{419, "The page expired, please try again.", "i-lucide-clock-alert"},
			{429, "You have made too many requests. Please wait and try again.", "i-lucide-timer"},
			{500, "Whoops, something went wrong on our end. Please try again.", "i-lucide-server-crash"},
			{502, "Whoops, something went wrong on our end. Please try again.", "i-lucide-server-crash"},
			{503, "Sorry, we are doing some maintenance. Please check back soon.", "i-lucide-construction"},
		} {
			request := httptest.NewRequest("POST", fmt.Sprintf("/_test/error/%d", test.status), nil)
			request.Header.Set("X-Inertia", "true")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			require.Equal(t, test.status, response.Code)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &payload))
			require.Equal(t, test.detail, payload["errorDetail"])
			require.Equal(t, test.icon, payload["errorIcon"])
			if test.status == 419 {
				require.Equal(t, "Page Expired - 419", payload["errorSummary"])
			}
		}
	})
	t.Run("partial status override inherits the default icon", func(t *testing.T) {
		original := config.Get("errors.statuses")
		config.Add("errors.statuses", map[string]any{
			"405": map[string]any{"detail": "Custom request message."},
		})
		defer config.Add("errors.statuses", original)
		for _, method := range []string{"GET", "POST"} {
			request := httptest.NewRequest(method, "/_test/error/405", nil)
			request.Header.Set("X-Inertia", "true")
			request.Header.Set("X-Inertia-Version", facades.Inertia().Version())
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			require.Equal(t, 405, response.Code)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &payload))
			if method == "GET" {
				require.Equal(t, "Custom request message.", payload["props"].(map[string]any)["detail"])
			} else {
				require.Equal(t, "Custom request message.", payload["errorDetail"])
				require.Equal(t, "i-lucide-circle-alert", payload["errorIcon"])
			}
		}
	})
	for _, path := range []string{"/api/missing", "/health/missing", "/api/_test/panic"} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		require.Contains(t, response.Header().Get("Content-Type"), "application/json")
		require.Empty(t, response.Header().Get("X-Inertia"))
		require.NotContains(t, response.Body.String(), "private database password")
	}
	t.Run("debug server errors keep the existing non-Inertia response", func(t *testing.T) {
		config.Add("app.debug", true)
		defer config.Add("app.debug", false)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest("GET", "/_test/panic", nil))
		require.Equal(t, 500, response.Code)
		require.Empty(t, response.Header().Get("X-Inertia"))
		require.NotContains(t, response.Body.String(), "private database password")
	})
}

func TestFlashToasts(t *testing.T) {
	router := facades.Route()
	router.Middleware(sessionmiddleware.StartSession()).Group(func(group route.Router) {
		group.Get("/_test/flash", func(ctx http.Context) http.Response {
			ctx.Request().Session().Flash("success_toast", "Saved successfully.")
			ctx.Request().Session().Flash("warn_toast", "Check your input.")
			ctx.Request().Session().Flash("message_toast", "An update is available.")
			return facades.Inertia().Redirect(ctx, "/")
		})
	})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest("GET", "/_test/flash", nil))
	require.Equal(t, 302, response.Code)
	cookies := response.Result().Cookies()
	for visit := 0; visit < 2; visit++ {
		request := httptest.NewRequest("GET", "/", nil)
		request.Header.Set("X-Inertia", "true")
		request.Header.Set("X-Inertia-Version", facades.Inertia().Version())
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)
		require.Equal(t, 200, response.Code)
		var page struct{ Props map[string]any }
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &page))
		if visit == 0 {
			require.Equal(t, "Saved successfully.", page.Props["flash"].(map[string]any)["success_toast"])
			require.Equal(t, "Check your input.", page.Props["flash"].(map[string]any)["warn_toast"])
			require.Equal(t, "An update is available.", page.Props["flash"].(map[string]any)["message_toast"])
		} else {
			require.Empty(t, page.Props["flash"])
		}
	}
}

func TestBuiltAssets(t *testing.T) {
	manifest, err := os.ReadFile("public/build/.vite/manifest.json")
	if os.IsNotExist(err) {
		t.Skip("run npm run build to check production asset routing")
	}
	require.NoError(t, err)
	var chunks map[string]struct {
		File string
		CSS  []string
	}
	require.NoError(t, json.Unmarshal(manifest, &chunks))
	entry := chunks["resources/js/app.ts"]
	require.NotEmpty(t, entry.File)
	require.NotEmpty(t, entry.CSS, "Nuxt UI and Tailwind styles must be built")
	for _, asset := range append([]string{entry.File}, entry.CSS...) {
		response := httptest.NewRecorder()
		facades.Route().ServeHTTP(response, httptest.NewRequest("GET", "/build/"+asset, nil))
		require.Equal(t, 200, response.Code)
		require.NotEmpty(t, response.Body.String())
		require.Empty(t, response.Header().Get("Set-Cookie"))
	}
	require.Equal(t, "", facades.Config().GetString("inertia.vite.dev_url"))
}

func TestPostgresIntegration(t *testing.T) {
	name := os.Getenv("GONERTIA_TEST_DATABASE")
	if name == "" {
		t.Skip("set GONERTIA_TEST_DATABASE to a separate gonertia_test_* database")
	}
	var databases []string
	require.NoError(t, facades.DB().Select(&databases, "SELECT current_database()"))
	require.Equal(t, []string{name}, databases)
	response := httptest.NewRecorder()
	facades.Route().ServeHTTP(response, httptest.NewRequest("GET", "/health/ready", nil))
	require.Equal(t, 200, response.Code)
	require.JSONEq(t, `{"status":"ok"}`, response.Body.String())

	migration := smokeMigration{}
	require.False(t, facades.Schema().HasTable("gonertia_migration_smoke"))
	facades.Schema().Register([]schema.Migration{migration})
	t.Cleanup(func() {
		require.NoError(t, facades.Schema().DropIfExists("gonertia_migration_smoke"))
	})
	require.NoError(t, facades.Artisan().Call("migrate"))
	require.True(t, facades.Schema().HasTable("gonertia_migration_smoke"))
	require.NoError(t, facades.DB().Statement("INSERT INTO gonertia_migration_smoke (id) VALUES (1)"))
	var ids []int
	require.NoError(t, facades.DB().Select(&ids, "SELECT id FROM gonertia_migration_smoke"))
	require.Equal(t, []int{1}, ids)
	require.NoError(t, facades.Artisan().Call("migrate:status"))
	require.NoError(t, facades.Artisan().Call("migrate:rollback --step=1"))
	require.False(t, facades.Schema().HasTable("gonertia_migration_smoke"))
}

// This fixture is registered only in the isolated database integration test.
type smokeMigration struct{}

func (smokeMigration) Signature() string { return "20261006000000_gonertia_migration_smoke" }
func (smokeMigration) Up() error {
	return facades.Schema().Sql("CREATE TABLE gonertia_migration_smoke (id INTEGER PRIMARY KEY)")
}
func (smokeMigration) Down() error { return facades.Schema().Drop("gonertia_migration_smoke") }
