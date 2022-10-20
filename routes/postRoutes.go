package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"net/http"
)

func PostRoutes(app *app.App) {
	// Get controller struct initialize
	postController := controllers.PostController{
		App: app,
	}

	http.HandleFunc("/register-handle", postController.Register)
	http.HandleFunc("/login-handle", postController.Login)
}
