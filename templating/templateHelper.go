package templating

import (
	"GoWeb/app"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
)

var templates = make(map[string]*template.Template) // This is only used here, does not need to be in app.App

func BuildPages(app *app.App) error {
	basePath := app.Config.Template.BaseName

	baseContent, err := app.Res.ReadFile(basePath)
	if err != nil {
		return fmt.Errorf("error reading base file: %w", err)
	}

	base, err := template.New(basePath).Parse(string(baseContent)) // Sets filepath as name and parses content
	if err != nil {
		return fmt.Errorf("error parsing base file: %w", err)
	}

	readFilesRecursively := func(fsys fs.FS, root string) ([]string, error) {
		var files []string
		err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking the path %q: %w", path, err)
			}
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		return files, err
	}

	// Get all file paths in the directory tree
	filePaths, err := readFilesRecursively(app.Res, app.Config.Template.ContentPath)
	if err != nil {
		return fmt.Errorf("error reading files recursively: %w", err)
	}

	for _, contentPath := range filePaths { // Create a new template base + content for each page
		content, err := app.Res.ReadFile(contentPath)
		if err != nil {
			return fmt.Errorf("error reading content file %s: %w", contentPath, err)
		}

		t, err := base.Clone()
		if err != nil {
			return fmt.Errorf("error cloning base template: %w", err)
		}

		_, err = t.Parse(string(content))
		if err != nil {
			return fmt.Errorf("error parsing content: %w", err)
		}

		templates[contentPath] = t
	}

	return nil
}

func RenderTemplate(w http.ResponseWriter, contentPath string, data any) {
	t, ok := templates[contentPath]
	if !ok {
		err := fmt.Errorf("template not found for path: %s", contentPath)
		slog.Error(err.Error())
		http.Error(w, "Template not found", 404)
		return
	}

	err := t.Execute(w, data) // Execute prebuilt template with dynamic data
	if err != nil {
		err = fmt.Errorf("error executing template: %w", err)
		slog.Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}
