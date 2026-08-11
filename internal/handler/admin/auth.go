package admin

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/pkg"
	"github.com/convee/goblog/internal/security"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var signinLimiter = security.NewLoginLimiter(5, 15*time.Minute)

func Login(c *artgo.Context) {
	data := make(map[string]interface{})
	pkg.AdminRender(data, c, "login")
}

func Logout(c *artgo.Context) {
	security.ClearSessionCookies(c.Writer, c.Req)
	c.Redirect(http.StatusFound, "/login")
}

func Signin(c *artgo.Context) {
	if !security.ValidateCSRFToken(c.Req, c.PostForm("_csrf")) {
		c.Error(http.StatusForbidden, "请求已失效，请刷新页面后重试")
		return
	}

	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	now := time.Now()
	limiterKey := signinLimiterKey(c.Req.RemoteAddr, email)
	if !signinLimiter.Allow(limiterKey, now) {
		c.Writer.Header().Set("Retry-After", "900")
		c.Error(http.StatusTooManyRequests, "登录尝试过多，请稍后重试")
		return
	}

	if email == "" || len(email) > 50 || password == "" || len(password) > 1024 {
		signinLimiter.RecordFailure(limiterKey, now)
		renderSigninError(c)
		return
	}
	user, err := daos.GetUser(email)
	passwordHash := user.Password
	if err != nil {
		// Keep failure work comparable for unknown and known accounts.
		passwordHash = "$2a$08$sWNweq0/ldMACfLWODHKg.6WKHYfjtyePQAAsVMryUySMSWDqSaKi"
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil || err != nil {
		signinLimiter.RecordFailure(limiterKey, now)
		renderSigninError(c)
		return
	}
	token, expires, err := security.NewSessionToken(user.Email, user.Password, time.Now())
	if err != nil {
		c.Error(http.StatusInternalServerError, "登录暂时不可用")
		return
	}
	signinLimiter.Reset(limiterKey)
	security.SetSessionCookie(c.Writer, c.Req, token, expires)
	c.Redirect(http.StatusFound, "/admin")
}

func signinLimiterKey(remoteAddr, email string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	return strings.ToLower(strings.TrimSpace(email)) + "|" + host
}

func renderSigninError(c *artgo.Context) {
	renderAdminError(c, http.StatusUnauthorized, "邮箱或密码不正确，请重试")
}
