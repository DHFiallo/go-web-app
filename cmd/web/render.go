package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type templateData struct {
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, t string, td *templateData) {
	var tmpl *template.Template

	// if we are using the template cache, try to get the template from our map, stored in the receiver
	if app.config.useCache {
		if templateFromMap, ok := app.templateMap[t]; ok {
			tmpl = templateFromMap
		}
	}

	if tmpl == nil {
		newTemplate, err := app.buildTemplateFromDisk(t)
		if err != nil {
			log.Println("Error building template: ", err)
			return
		}
		log.Println("Building template from disk")
		tmpl = newTemplate
	}

	if td == nil {
		td = &templateData{}
	}

	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		log.Println("Error executing template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) buildTemplateFromDisk(t string) (*template.Template, error) {
	templateSlice := []string{
		"./cmd/templates/base.layout.gohtml",
		"./cmd/templates/partials/header.partial.gohtml",
		"./cmd/templates/partials/footer.partial.gohtml",
		fmt.Sprintf("./cmd/templates/%s", t),
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", dir)

	tmpl, err := template.ParseFiles(templateSlice...)

	if err != nil {
		return nil, err
	}

	app.templateMap[t] = tmpl

	return tmpl, nil
}
