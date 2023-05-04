package controllers

import (
	"GoWeb/app"
	"GoWeb/models"
	"GoWeb/security"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

func (postController *PostController) FileUpload(w http.ResponseWriter, r *http.Request) {
	max := postController.App.Config.Upload.MaxSize
	err := r.ParseMultipartForm(max)
	if err != nil {
		return
	}

	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader, so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	if handler.Size > max {
		log.Println("User tried uploading a file which is too large.")
		http.Redirect(w, r, "/", http.StatusRequestHeaderFieldsTooLarge)
		return
	}

	// Create a temporary file within upload directory
	tempFile, err := os.Create(postController.App.Config.Upload.BaseName + handler.Filename)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusNotAcceptable)
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(tempFile)

	// read all the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
