package middleware

import (
	"GoWeb/security"
	"log"
	"net/http"
)

// Csrf validates the CSRF token and returns the handler function if it succeded
func Csrf(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := security.VerifyCsrfToken(r)
		if err != nil {
			log.Println("Error verifying csrf token")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		f(w, r)
	}
}
