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

func (g *Get) ShowHome(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken       string
		IsAuthenticated bool
		Test            string
	}

	CsrfToken, err := security.GenerateCsrfToken(w, r)
	if err != nil {
		return
	}

	isAuthenticated := true
	user, err := models.CurrentUser(g.App, r)
	if err != nil || user.Id == 0 {
		isAuthenticated = false
	}

	data := dataStruct{
		CsrfToken:       CsrfToken,
		Test:            "Hello World!",
		IsAuthenticated: isAuthenticated,
	}

	templating.RenderTemplate(w, "templates/pages/home.html", data)
}

func (g *Get) ShowRegister(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken       string
		IsAuthenticated bool
	}

	CsrfToken, err := security.GenerateCsrfToken(w, r)
	if err != nil {
		return
	}

	isAuthenticated := true
	user, err := models.CurrentUser(g.App, r)
	if err != nil || user.Id == 0 {
		isAuthenticated = false
	}

	data := dataStruct{
		CsrfToken:       CsrfToken,
		IsAuthenticated: isAuthenticated,
	}

	templating.RenderTemplate(w, "templates/pages/register.html", data)
}

func (g *Get) ShowLogin(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken       string
		IsAuthenticated bool
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
