package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"io/fs"
	"log/slog"
	"net/http"
)

// Get defines all project get routes
func Get(app *app.App) {
	// Get controller struct initialize
	getController := controllers.Get{
		App: app,
	}

	// Serve static files
	staticFS, err := fs.Sub(app.Res, "static")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	staticHandler := http.FileServer(http.FS(staticFS))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	slog.Info("serving static files from embedded file system /static")

	// Pages
	http.HandleFunc("GET /", getController.ShowHome)
	http.HandleFunc("GET /login", getController.ShowLogin)
	http.HandleFunc("GET /register", getController.ShowRegister)
}
