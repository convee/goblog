package routers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/security"
)

func TestLoginPageHasSecurityControls(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file")
	}
	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	conf.Conf = &conf.Config{App: conf.AppConfig{Cdn: "/static"}}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://example.com/login", nil)
	InitRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /login status = %d, want 200", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `name="_csrf"`) {
		t.Fatal("login form is missing its CSRF field")
	}
	if recorder.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("security headers were not applied")
	}
	foundCSRF := false
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == security.CSRFCookieName {
			foundCSRF = cookie.HttpOnly && cookie.Secure && cookie.SameSite == http.SameSiteStrictMode
		}
	}
	if !foundCSRF {
		t.Fatal("secure CSRF cookie was not issued")
	}
}
