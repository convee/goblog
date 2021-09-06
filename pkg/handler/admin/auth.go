package admin

import (
	"blog/pkg/model"
	"blog/pkg/mysql"
	"blog/pkg/view"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	view.AdminRender(data, w, "login")
}

func Register(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	view.AdminRender(data, w, "register")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "email",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	repassword := r.FormValue("repassword")
	if email == "" || password == "" || repassword == "" || password != repassword {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	user := model.User{
		Email:    email,
		Password: string(hashPassword),
	}
	if _, err := mysql.AddUser(user); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cookie := &http.Cookie{
		Name:  "email",
		Value: email,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := mysql.GetUser(email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cookie := &http.Cookie{
		Name:  "email",
		Value: email,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin", http.StatusFound)
}
