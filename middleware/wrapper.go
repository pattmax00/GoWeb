package middleware

import "net/http"

// ProcessMiddleware is a wrapper function for the http.HandleFunc function
// that takes the function you want to execute (f) and the middleware you want
// to execute (m) this should be used when processing multiple groups of middleware at a time
func ProcessMiddleware(f func(w http.ResponseWriter, r *http.Request), m []func()) func(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range m {
		middleware()
	}

	return f
}
