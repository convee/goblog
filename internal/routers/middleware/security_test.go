package middleware

import (
	"github.com/convee/goblog/conf"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/security"
)

func TestContentSecurityPolicyUsesConfiguredCDNOrigin(t *testing.T) {
	policy := contentSecurityPolicy("https://cdn.example.com/static")
	if !strings.Contains(policy, "script-src 'self' https://cdn.example.com") {
		t.Fatalf("script source missing configured CDN: %s", policy)
	}
	if strings.Contains(policy, "script-src 'self' 'unsafe-inline'") {
		t.Fatalf("inline scripts should not be allowed: %s", policy)
	}
}

func TestSecurityHeadersDoNotTrustForwardedProtoByDefault(t *testing.T) {
	previous := conf.Conf
	conf.Conf = &conf.Config{App: conf.AppConfig{Cdn: "/static"}}
	t.Cleanup(func() { conf.Conf = previous })
	if err := security.ConfigureTrustedProxies(nil); err != nil {
		t.Fatal(err)
	}

	router := artgo.New()
	router.Use(SecurityHeaders())
	router.GET("/", func(c *artgo.Context) { c.Status(http.StatusOK) })
	request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("untrusted forwarded proto enabled HSTS")
	}
}

func TestRequireAuthDoesNotInvokeRouteWithoutSession(t *testing.T) {
	invoked := false
	router := artgo.New()
	router.GET("/admin", RequireAuth(func(c *artgo.Context) {
		invoked = true
		c.Status(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/admin", nil))
	if invoked {
		t.Fatal("protected route was invoked without a session")
	}
	if recorder.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusFound)
	}
}

func TestRequireCSRFStopsAndAllowsRoute(t *testing.T) {
	invocations := 0
	router := artgo.New()
	router.POST("/write", RequireCSRF(func(c *artgo.Context) {
		invocations++
		c.Status(http.StatusNoContent)
	}))

	rejected := httptest.NewRecorder()
	router.ServeHTTP(rejected, httptest.NewRequest(http.MethodPost, "/write", nil))
	if rejected.Code != http.StatusForbidden || invocations != 0 {
		t.Fatalf("missing token: status=%d invocations=%d", rejected.Code, invocations)
	}

	tokenRecorder := httptest.NewRecorder()
	tokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	token, err := security.EnsureCSRFToken(tokenRecorder, tokenRequest)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{"_csrf": {token}}
	request := httptest.NewRequest(http.MethodPost, "/write", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.AddCookie(tokenRecorder.Result().Cookies()[0])
	accepted := httptest.NewRecorder()
	router.ServeHTTP(accepted, request)
	if accepted.Code != http.StatusNoContent || invocations != 1 {
		t.Fatalf("valid token: status=%d invocations=%d", accepted.Code, invocations)
	}
}
