package middleware

import "net/http"

func AuthWithCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("email"); err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}
		next.ServeHTTP(w, r)
	})
}
