package middleware

import "net/http"

// ProcessGroup is a wrapper function for the http.HandleFunc function
// that takes the function you want to execute (f) and the middleware you want
// to execute (m) this should be used when processing multiple groups of middleware at a time
func ProcessGroup(f func(w http.ResponseWriter, r *http.Request), m []MiddlewareFunc) func(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range m {
		_ = middleware(f)
	}

	return f
}
