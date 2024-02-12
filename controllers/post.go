package controllers

import (
	"GoWeb/app"
	"GoWeb/models"
	"log/slog"
	"net/http"
	"time"
)

// Post is a wrapper struct for the App struct
type Post struct {
	App *app.App
}

func (p *Post) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "on"

	if username == "" || password == "" {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}

	_, err := models.AuthenticateUser(p.App, w, username, password, remember)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (p *Post) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	createdAt := time.Now()
	updatedAt := time.Now()

	if username == "" || password == "" {
		http.Redirect(w, r, "/register", http.StatusUnauthorized)
	}

	_, err := models.CreateUser(p.App, username, password, createdAt, updatedAt)
	if err != nil {
		// TODO: if err == bcrypt.ErrPasswordTooLong display error to user, this will require a flash message system with cookies
		slog.Error("error creating user: " + err.Error())
		http.Redirect(w, r, "/register", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (p *Post) Logout(w http.ResponseWriter, r *http.Request) {
	models.LogoutUser(p.App, w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
