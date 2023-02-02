package middleware

import "net/http"

func SecureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "HttpOnly")
		next.ServeHTTP(w, r)
	})
}
