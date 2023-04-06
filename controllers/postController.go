package controllers

import (
	"GoWeb/app"
	"GoWeb/models"
	"GoWeb/security"
	"log"
	"net/http"
	"time"
)

// PostController is a wrapper struct for the App struct
type PostController struct {
	App *app.App
}

func (postController *PostController) Login(w http.ResponseWriter, r *http.Request) {
	// Validate csrf token
	_, err := security.VerifyCsrfToken(r)
	if err != nil {
		log.Println("Error verifying csrf token")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "on"

	if username == "" || password == "" {
		log.Println("Tried to login user with empty username or password")
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	_, err = models.AuthenticateUser(postController.App, w, username, password, remember)
	if err != nil {
		log.Println("Error authenticating user")
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (postController *PostController) Register(w http.ResponseWriter, r *http.Request) {
	// Validate csrf token
	_, err := security.VerifyCsrfToken(r)
	if err != nil {
		log.Println("Error verifying csrf token")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	createdAt := time.Now()
	updatedAt := time.Now()

	if username == "" || password == "" {
		log.Println("Tried to create user with empty username or password")
		http.Redirect(w, r, "/register", http.StatusFound)
	}

	_, err = models.CreateUser(postController.App, username, password, createdAt, updatedAt)
	if err != nil {
		log.Println("Error creating user")
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
