package templating

import (
	"GoWeb/app"
	"html/template"
	"log"
	"net/http"
)

// RenderTemplate renders and serves a template from the embedded filesystem optionally with given data
func RenderTemplate(app *app.App, w http.ResponseWriter, contentPath string, data any) {
	templatePath := app.Config.Template.BaseName

	templateContent, err := app.Res.ReadFile(templatePath)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.New(templatePath).Parse(string(templateContent))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	content, err := app.Res.ReadFile(contentPath)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	t, err = t.Parse(string(content))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
