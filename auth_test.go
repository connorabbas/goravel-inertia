package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"goravel/app/facades"
	"goravel/app/models"
)

func TestAuthenticationJourney(t *testing.T) {
	if os.Getenv("GONERTIA_TEST_DATABASE") == "" {
		t.Skip("requires a separate gonertia_test_* database")
	}
	cookies := map[string]*http.Cookie{}
	visit := func(method, path string, fields url.Values, csrf bool) *httptest.ResponseRecorder {
		t.Helper()
		var body *strings.Reader
		if fields == nil {
			body = strings.NewReader("")
		} else {
			body = strings.NewReader(fields.Encode())
		}
		request := httptest.NewRequest(method, path, body)
		request.Host = "goravel.localhost"
		request.Header.Set("Referer", "http://goravel.localhost"+path)
		if fields != nil {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
		if csrf && cookies["XSRF-TOKEN"] != nil {
			request.Header.Set("X-CSRF-TOKEN", cookies["XSRF-TOKEN"].Value)
		}
		response := httptest.NewRecorder()
		facades.Route().ServeHTTP(response, request)
		for _, cookie := range response.Result().Cookies() {
			cookies[cookie.Name] = cookie
		}
		return response
	}
	response := visit("GET", "/register", nil, false)
	require.Equal(t, 200, response.Code)
	require.NotNil(t, cookies["XSRF-TOKEN"])
	response = visit("POST", "/register", url.Values{"name": {"Test User"}, "email": {"auth-journey@example.test"}, "password": {"Secretpassword123"}, "password_confirmation": {"Secretpassword123"}}, false)
	require.Equal(t, 419, response.Code, "registration requires a CSRF token")
	for _, invalid := range []struct {
		fields url.Values
		field  string
	}{
		{url.Values{"name": {" "}, "email": {"invalid"}, "password": {"short"}, "password_confirmation": {"different"}}, "email"},
		{url.Values{"name": {"Test User"}, "email": {"auth-journey@example.test"}, "password": {"Secretpassword123"}, "password_confirmation": {"different"}}, "password"},
	} {
		response = visit("POST", "/register", invalid.fields, true)
		require.Equal(t, 302, response.Code)
		response = visit("GET", "/register", nil, false)
		props := inertiaProps(t, response)
		require.NotEmpty(t, props["errors"].(map[string]any)[invalid.field])
		require.NotContains(t, response.Body.String(), "Secretpassword123")
	}
	response = visit("POST", "/register", url.Values{"name": {" Test User "}, "email": {" AUTH-JOURNEY@EXAMPLE.TEST "}, "password": {"Secretpassword123"}, "password_confirmation": {"Secretpassword123"}}, true)
	require.Equal(t, 302, response.Code, response.Body.String())
	require.Equal(t, "/verify-email", response.Header().Get("Location"))
	var user models.User
	require.NoError(t, facades.Orm().Query().Where("email", "auth-journey@example.test").First(&user))
	require.NotZero(t, user.ID)
	require.Equal(t, "Test User", user.Name)
	t.Cleanup(func() { facades.Orm().Query().Delete(&user) })
	require.NotEqual(t, "Secretpassword123", user.Password)
	response = visit("GET", "/dashboard", nil, false)
	require.Equal(t, 302, response.Code)
	require.Equal(t, "/verify-email", response.Header().Get("Location"))
	key := facades.Config().GetString("app.key")
	if len(key) < 32 {
		original := key
		key = "test-only-32-character-signing-key!!!"
		facades.Config().Add("app.key", key)
		t.Cleanup(func() { facades.Config().Add("app.key", original) })
	}
	expires := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(strconv.FormatUint(uint64(user.ID), 10) + ":" + user.Email + ":" + expires))
	verification := "/verify-email/" + strconv.FormatUint(uint64(user.ID), 10) + "/" + hex.EncodeToString(mac.Sum(nil)) + "?expires=" + expires
	response = visit("GET", strings.Replace(verification, "?", "0?", 1), nil, false)
	require.Equal(t, 403, response.Code, "tampered verification links are rejected")
	response = visit("GET", verification, nil, false)
	require.Equal(t, 302, response.Code)
	require.Equal(t, "/dashboard", response.Header().Get("Location"))
	response = visit("GET", "/dashboard", nil, false)
	require.Equal(t, 200, response.Code)
	require.NotContains(t, response.Body.String(), user.Password)
	var page map[string]any
	_, body, ok := strings.Cut(response.Body.String(), `<script data-page="app" type="application/json">`)
	require.True(t, ok)
	body, _, _ = strings.Cut(body, "</script>")
	require.NoError(t, json.Unmarshal([]byte(body), &page))
	auth := page["props"].(map[string]any)["auth"].(map[string]any)
	require.Equal(t, "auth-journey@example.test", auth["user"].(map[string]any)["email"])
	response = visit("PATCH", "/settings/profile", url.Values{"name": {strings.Repeat("界", 200)}, "email": {" AUTH-JOURNEY@EXAMPLE.TEST "}}, true)
	require.Equal(t, 303, response.Code)
	require.NoError(t, facades.Orm().Query().Find(&user, user.ID))
	require.Equal(t, strings.Repeat("界", 200), user.Name, "max validates characters, not UTF-8 bytes")
	require.NotNil(t, user.EmailVerifiedAt, "normalizing the unchanged email preserves verification")
	response = visit("PUT", "/settings/password", url.Values{"current_password": {"wrong"}, "password": {"Newsecretpassword123"}, "password_confirmation": {"Newsecretpassword123"}}, true)
	require.Equal(t, 303, response.Code)
	response = visit("PUT", "/settings/password", url.Values{"current_password": {"Secretpassword123"}, "password": {"Newsecretpassword123"}, "password_confirmation": {"Newsecretpassword123"}}, true)
	require.Equal(t, 303, response.Code)
	require.NoError(t, facades.Orm().Query().Where("email", "auth-journey@example.test").First(&user))
	require.True(t, facades.Hash().Check("Newsecretpassword123", user.Password))
	response = visit("POST", "/logout", nil, true)
	require.Equal(t, 303, response.Code)
	response = visit("GET", "/dashboard", nil, false)
	require.Equal(t, 302, response.Code)
	require.Equal(t, "/login", response.Header().Get("Location"))
	response = visit("GET", "/login", nil, false)
	require.Equal(t, 200, response.Code)
	response = visit("POST", "/register", url.Values{"name": {"Another User"}, "email": {" AUTH-JOURNEY@EXAMPLE.TEST "}, "password": {"Secretpassword123"}, "password_confirmation": {"Secretpassword123"}}, true)
	require.Equal(t, 302, response.Code)
	response = visit("GET", "/register", nil, false)
	require.NotEmpty(t, inertiaProps(t, response)["errors"].(map[string]any)["email"], "Goravel unique rule rejects a normalized duplicate")
	response = visit("POST", "/login", url.Values{"email": {"auth-journey@example.test"}, "password": {"Secretpassword123"}}, true)
	require.Equal(t, 302, response.Code)
	response = visit("GET", "/login", nil, false)
	response = visit("POST", "/login", url.Values{"email": {"auth-journey@example.test"}, "password": {"Newsecretpassword123"}}, true)
	require.Equal(t, 302, response.Code)
	require.Equal(t, "/dashboard", response.Header().Get("Location"))
	response = visit("GET", "/dashboard", nil, false)
	require.Equal(t, 200, response.Code)
	response = visit("POST", "/logout", nil, true)
	require.Equal(t, 303, response.Code)
	response = visit("GET", "/login", nil, false)
	require.Equal(t, 200, response.Code)
	response = visit("POST", "/login", url.Values{"email": {"auth-journey@example.test"}, "password": {"Newsecretpassword123"}, "remember": {"true"}}, true)
	require.Equal(t, 302, response.Code)
	response = visit("GET", "/dashboard", nil, false)
	require.Equal(t, 200, response.Code)
	require.NotNil(t, cookies["remember_user"])
	delete(cookies, facades.Config().GetString("session.cookie"))
	response = visit("GET", "/dashboard", nil, false)
	require.Equal(t, 200, response.Code, "remember cookie restores a session")
	response = visit("POST", "/logout", nil, true)
	require.Equal(t, 303, response.Code)
	delete(cookies, facades.Config().GetString("session.cookie"))
	response = visit("GET", "/dashboard", nil, false)
	require.Equal(t, 302, response.Code, "logout revokes the remember cookie")
	token := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	require.NoError(t, facades.Orm().Query().Create(&models.PasswordResetToken{Email: user.Email, TokenHash: hashTokenForTest(token + ":" + strconv.FormatUint(uint64(user.ID), 10)), ExpiresAt: time.Now().Add(time.Hour)}))
	response = visit("GET", "/reset-password/"+token+"?email=auth-journey@example.test", nil, false)
	require.Equal(t, 200, response.Code)
	reset := url.Values{"email": {user.Email}, "token": {token}, "password": {"Resettedpassword123"}, "password_confirmation": {"Resettedpassword123"}}
	response = visit("POST", "/reset-password", reset, true)
	require.Equal(t, 302, response.Code)
	require.NoError(t, facades.Orm().Query().Where("email", user.Email).First(&user))
	require.True(t, facades.Hash().Check("Resettedpassword123", user.Password))
	response = visit("POST", "/reset-password", reset, true)
	require.Equal(t, 302, response.Code, "reset tokens must be single use")
	require.NotEqual(t, "/login", response.Header().Get("Location"))
	require.NoError(t, facades.Orm().Query().Where("email", user.Email).First(&user))
	require.True(t, facades.Hash().Check("Resettedpassword123", user.Password))
}

func inertiaProps(t *testing.T, response *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	_, body, found := strings.Cut(response.Body.String(), `<script data-page="app" type="application/json">`)
	require.True(t, found)
	body, _, _ = strings.Cut(body, "</script>")
	var page struct{ Props map[string]any }
	require.NoError(t, json.Unmarshal([]byte(body), &page))
	return page.Props
}

func hashTokenForTest(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
