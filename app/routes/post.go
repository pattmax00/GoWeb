package routes

import (
	"GoWeb/app"
	"GoWeb/app/controllers"
	"GoWeb/app/middleware"
	"net/http"
)

// Post defines all project post routes
func Post(app *app.Deps) {
	// Post controller struct initialize
	postController := controllers.Post{
		App: app,
	}

	// User authentication
	http.HandleFunc("/register-handle", middleware.Csrf(postController.Register))
	http.HandleFunc("/login-handle", middleware.Csrf(postController.Login))
}
