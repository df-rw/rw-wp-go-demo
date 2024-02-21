package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Application struct {
	tpl *template.Template
}

func (app *Application) render(w http.ResponseWriter, pageName string, pageData any,
	statusCode int) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	app.tpl.ExecuteTemplate(w, pageName, pageData)
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, "home", nil, http.StatusOK)
}

func (app *Application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, "about", nil, http.StatusOK)
}

func (app *Application) form(w http.ResponseWriter, r *http.Request) {
	app.render(w, "form", nil, http.StatusOK)
}

func (app *Application) processForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	searchTerm := r.FormValue("search-term")

	var message string
	var results []string

	if searchTerm == "" {
		message = "You didn't specify a search term"
	} else {
		fruits := []string{
			"apple", "banana", "pear", "grapefruit", "kiwi", "pineapple",
		}

		for _, v := range fruits {
			if strings.Contains(v, searchTerm) {
				results = append(results, v)
			}
		}

		if len(results) == 0 {
			message = "No results returned!"
		} else {
			message = fmt.Sprintf("Got %d results!", len(results))
		}
	}

	pageData := map[string]any{
		"Message": message,
		"Results": results,
	}

	app.render(w, "fruit-results", pageData, http.StatusOK)
}

func main() {
	tpl := template.Must(template.ParseGlob("*.tmpl"))
	app := &Application{tpl}

	http.HandleFunc("/", app.home)
	http.HandleFunc("/about", app.about)
	http.HandleFunc("/form", app.form)
	http.HandleFunc("/process-form", app.processForm)

	fmt.Println("listening on 4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
