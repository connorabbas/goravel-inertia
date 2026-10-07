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
	"github.com/goravel/framework/foundation"
	frameworkmock "github.com/goravel/framework/testing/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goravel/app/facades"
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
			require.Contains(t, response.Body.String(), `<div id="app"></div>`)
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

func TestBuiltAssets(t *testing.T) {
	manifest, err := os.ReadFile("public/build/.vite/manifest.json")
	if os.IsNotExist(err) {
		t.Skip("run npm run build to check production asset routing")
	}
	require.NoError(t, err)
	var chunks map[string]struct{ File string }
	require.NoError(t, json.Unmarshal(manifest, &chunks))
	entry := chunks["resources/js/app.ts"].File
	require.NotEmpty(t, entry)
	response := httptest.NewRecorder()
	facades.Route().ServeHTTP(response, httptest.NewRequest("GET", "/build/"+entry, nil))
	require.Equal(t, 200, response.Code)
	require.NotEmpty(t, response.Body.String())
	require.Empty(t, response.Header().Get("Set-Cookie"))
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
