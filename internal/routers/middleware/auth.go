package middleware

import (
	"github.com/convee/artgo"
	"net/http"
)

func AuthWithCookie() artgo.HandlerFunc {
	return func(c *artgo.Context) {
		if cookie, err := c.Req.Cookie("email"); err != nil || cookie.Value == "" {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.Next()
	}
}
