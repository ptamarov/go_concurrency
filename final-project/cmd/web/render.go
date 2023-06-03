package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

// something to store data in that will be passed to templates

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	// User *data.User <- later
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	var templateSlice []string

	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, t))

	templateSlice = append(templateSlice, partials...)

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)

	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, app.addDefaultData(td, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Config) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	// make sure that this info is removed from session as soon as it is read

	if app.isAuthenticated(r) {
		td.Authenticated = true
		// TODO: get more user information
	}
	td.Now = time.Now()

	return td
}

func (app *Config) isAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
