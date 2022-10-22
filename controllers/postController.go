package controllers

import (
	"GoWeb/app"
	"GoWeb/database/models"
	"log"
	"net/http"
	"time"
)

// PostController is a wrapper struct for the App struct
type PostController struct {
	App *app.App
}

func (postController *PostController) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	createdAt := time.Now()
	updatedAt := time.Now()

	if username == "" || password == "" {
		log.Println("Tried to create user with empty username or password")
		http.Redirect(w, r, "/register", http.StatusFound)
	}

	_, err := models.CreateUser(postController.App, username, password, createdAt, updatedAt)
	if err != nil {
		log.Println("Error creating user")
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (postController *PostController) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		log.Println("Tried to create user with empty username or password")
		http.Redirect(w, r, "/register", http.StatusFound)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
