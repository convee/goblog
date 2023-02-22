package admin

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/pkg"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(c *artgo.Context) {
	data := make(map[string]interface{})
	pkg.AdminRender(data, c, "login")
}

func Register(c *artgo.Context) {
	data := make(map[string]interface{})
	pkg.AdminRender(data, c, "register")
}

func Logout(c *artgo.Context) {
	c.SetCookie(&http.Cookie{
		Name:   "email",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	c.Redirect(http.StatusFound, "/login")
}

func Signup(c *artgo.Context) {
	c.Redirect(http.StatusFound, "/admin")
}

func Signin(c *artgo.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.Status(http.StatusInternalServerError)
		return
	}
	user := daos.GetUser(email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data := make(map[string]interface{})
		data["msg"] = "密码不正确，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}
	cookie := &http.Cookie{
		Name:  "email",
		Value: email,
		Path:  "/",
	}
	c.SetCookie(cookie)
	c.Redirect(http.StatusFound, "/admin")
	return
}
