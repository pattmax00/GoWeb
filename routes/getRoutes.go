package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"io/fs"
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
	staticFS, err := fs.Sub(app.Res, "static")
	if err != nil {
		log.Println(err)
		return
	}
	staticHandler := http.FileServer(http.FS(staticFS))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	log.Println("Serving static files from embedded file system /static")

	// Pages
	http.HandleFunc("/", getController.ShowHome)
	http.HandleFunc("/login", getController.ShowLogin)
	http.HandleFunc("/register", getController.ShowRegister)
	http.HandleFunc("/logout", getController.Logout)
}
