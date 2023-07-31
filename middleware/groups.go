package middleware

import "net/http"

type MiddlewareFunc func(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
