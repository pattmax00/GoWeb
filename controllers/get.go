package controllers

import (
	"GoWeb/app"
	"GoWeb/models"
	"GoWeb/security"
	"GoWeb/templating"
	"net/http"
)

// Get is a wrapper struct for the App struct
type Get struct {
	App *app.App
}

func (g *Get) ShowHome(w http.ResponseWriter, _ *http.Request) {
	type dataStruct struct {
		Test string
	}

	data := dataStruct{
		Test: "Hello World!",
	}

	templating.RenderTemplate(w, "templates/pages/home.html", data)
}

func (g *Get) ShowRegister(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken string
	}

	CsrfToken, err := security.GenerateCsrfToken(w, r)
	if err != nil {
		return
	}

	data := dataStruct{
		CsrfToken: CsrfToken,
	}

	templating.RenderTemplate(w, "templates/pages/register.html", data)
}

func (g *Get) ShowLogin(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken string
	}

	CsrfToken, err := security.GenerateCsrfToken(w, r)
	if err != nil {
		return
	}

	data := dataStruct{
		CsrfToken: CsrfToken,
	}

	templating.RenderTemplate(w, "templates/pages/login.html", data)
}

func (g *Get) Logout(w http.ResponseWriter, r *http.Request) {
	models.LogoutUser(g.App, w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
