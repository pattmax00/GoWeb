package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"GoWeb/middleware"
	"net/http"
)

// Post defines all project post routes
func Post(app *app.App) {
	// Post controller struct initialize
	postController := controllers.Post{
		App: app,
	}

	// User authentication
	http.HandleFunc("POST /register-handle", middleware.Csrf(postController.Register))
	http.HandleFunc("POST /login-handle", middleware.Csrf(postController.Login))
	http.HandleFunc("POST /logout", middleware.Csrf(postController.Logout))
}
