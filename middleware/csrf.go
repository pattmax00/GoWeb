package middleware

import (
	"GoWeb/security"
	"log"
	"net/http"
)

// CsrfMiddleware validates the CSRF token and returns the handler function if it succeded
func CsrfMiddleware(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify csrf token
		_, err := security.VerifyCsrfToken(r)
		if err != nil {
			log.Println("Error verifying csrf token")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		f(w, r)
	}
}
