package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"log"
	"net/http"
)

// GetRoutes defines all project get routes
func GetRoutes(app *app.App) {
	// Get controller struct initialize
	getController := controllers.GetController{
		App: app,
	}

	// Serve static files
	http.Handle("/file/", http.FileServer(http.Dir("./static")))
	log.Println("Serving static files from: ./static")

	// Pages
	http.HandleFunc("/", getController.ShowHome)
	http.HandleFunc("/register", getController.ShowRegister)
}
