package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("template/edit.gohtml", "template/view.gohtml"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	http.HandleFunc("/save/", makeHandler(saveHandler)) // url: "/save/foo" will create foo file
	http.HandleFunc("/edit/", makeHandler(editHandler)) // url: "edit/foo" allows user to edit content in file
	http.HandleFunc("/view/", makeHandler(viewHandler)) // url: "/view"
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (p *Page) saveAsFile() error {
	fileName := p.Title + ".txt"

	return ioutil.WriteFile(fileName, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	fileName := title + ".txt"
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate("view", w, page)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate("edit", w, page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	err := page.saveAsFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(tmpl string, w http.ResponseWriter, page *Page) {
	fullTmpl := tmpl + ".gohtml"
	err := templates.ExecuteTemplate(w, fullTmpl, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		match := validPath.FindStringSubmatch(r.URL.Path)
		if match == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, match[2])
	}
}
