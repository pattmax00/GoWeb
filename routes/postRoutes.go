package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"GoWeb/middleware"
	"net/http"
)

// PostRoutes defines all project post routes
func PostRoutes(app *app.App) {
	// Post controller struct initialize
	postController := controllers.PostController{
		App: app,
	}

	// User authentication
	http.HandleFunc("/register-handle", middleware.CsrfMiddleware(postController.Register))
	http.HandleFunc("/login-handle", middleware.CsrfMiddleware(postController.Login))
}
