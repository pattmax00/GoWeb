package controllers

import (
	"GoWeb/app"
	"GoWeb/database/models"
	"GoWeb/templating"
	"net/http"
)

// GetController is a wrapper struct for the App struct
type GetController struct {
	App *app.App
}

func (getController *GetController) ShowHome(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		Test string
	}

	data := dataStruct{
		Test: "Hello World!",
	}

	templating.RenderTemplate(getController.App, w, "templates/pages/home.html", data)
}

func (getController *GetController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	templating.RenderTemplate(getController.App, w, "templates/pages/register.html", nil)
}

func (getController *GetController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	templating.RenderTemplate(getController.App, w, "templates/pages/login.html", nil)
}

func (getController *GetController) Logout(w http.ResponseWriter, r *http.Request) {
	models.LogoutUser(getController.App, w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
