package middleware

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/security"
	"net/http"
	"net/url"
	"strings"
)

const maxRequestBodyBytes = 2 << 20

func SecurityHeaders() artgo.HandlerFunc {
	return func(c *artgo.Context) {
		header := c.Writer.Header()
		header.Set("Content-Security-Policy", contentSecurityPolicy(conf.Conf.App.Cdn))
		header.Set("Cross-Origin-Opener-Policy", "same-origin")
		header.Set("Permissions-Policy", "camera=(), geolocation=(), microphone=()")
		header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		header.Set("X-Content-Type-Options", "nosniff")
		header.Set("X-Frame-Options", "DENY")
		header.Set("X-Permitted-Cross-Domain-Policies", "none")
		if c.Req.URL.Path == "/login" || c.Req.URL.Path == "/signin" || strings.HasPrefix(c.Req.URL.Path, "/admin") {
			header.Set("Cache-Control", "no-store")
		} else if strings.HasPrefix(c.Req.URL.Path, "/static/") {
			header.Set("Cache-Control", "public, max-age=3600")
		}
		if security.RequestIsHTTPS(c.Req) {
			header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

func contentSecurityPolicy(cdn string) string {
	source := ""
	if parsed, err := url.Parse(cdn); err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != "" {
		source = " " + parsed.Scheme + "://" + parsed.Host
	}
	return "default-src 'self'; " +
		"base-uri 'self'; " +
		"connect-src 'self'" + source + "; " +
		"font-src 'self' data:" + source + "; " +
		"form-action 'self'; " +
		"frame-ancestors 'none'; " +
		"img-src 'self' data: http: https:; " +
		"object-src 'none'; " +
		"script-src 'self'" + source + "; " +
		"style-src 'self' 'unsafe-inline'" + source
}

func LimitRequestBody() artgo.HandlerFunc {
	return func(c *artgo.Context) {
		if c.Req.Body != nil && c.Req.Method != http.MethodGet && c.Req.Method != http.MethodHead {
			c.Req.Body = http.MaxBytesReader(c.Writer, c.Req.Body, maxRequestBodyBytes)
		}
		c.Next()
	}
}

func RequireCSRF(next artgo.HandlerFunc) artgo.HandlerFunc {
	return func(c *artgo.Context) {
		if c.Req.Method == http.MethodPost && !security.ValidateCSRFToken(c.Req, c.PostForm("_csrf")) {
			c.Error(http.StatusForbidden, "请求已失效，请刷新页面后重试")
			return
		}
		next(c)
	}
}

func AdminOnly(next artgo.HandlerFunc) artgo.HandlerFunc {
	return RequireAuth(RequireCSRF(next))
}
