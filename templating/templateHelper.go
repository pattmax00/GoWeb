package templating

import (
	"GoWeb/app"
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(app *app.App, w http.ResponseWriter, contentPath string, data any) {
	templatePath := app.Config.Template.BaseName

	t, _ := template.ParseFiles(templatePath, contentPath)
	err := t.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
