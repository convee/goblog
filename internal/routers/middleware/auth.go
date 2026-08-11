package middleware

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/security"
	"net/http"
	"time"
)

// RequireAuth wraps a route handler instead of relying on artgo's middleware
// short-circuit behavior. artgo v1.1.0 continues its middleware loop after a
// middleware returns, so authorization must own whether the route is invoked.
func RequireAuth(next artgo.HandlerFunc) artgo.HandlerFunc {
	return func(c *artgo.Context) {
		cookie, err := c.Req.Cookie(security.SessionCookieName)
		if err != nil || cookie.Value == "" {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		email, err := security.SessionEmail(cookie.Value)
		if err != nil {
			security.ClearSessionCookies(c.Writer, c.Req)
			c.Redirect(http.StatusFound, "/login")
			return
		}
		user, err := daos.GetUser(email)
		if err != nil {
			security.ClearSessionCookies(c.Writer, c.Req)
			c.Redirect(http.StatusFound, "/login")
			return
		}
		if _, err := security.ValidateSessionToken(cookie.Value, user.Password, time.Now()); err != nil {
			security.ClearSessionCookies(c.Writer, c.Req)
			c.Redirect(http.StatusFound, "/login")
			return
		}
		next(c)
	}
}
