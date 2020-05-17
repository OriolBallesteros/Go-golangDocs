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
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
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
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate("view", w, page)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate("edit", w, page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }
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
	fullTmpl := "template/" + tmpl + ".gohtml"
	err := templates.ExecuteTemplate(w, fullTmpl, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
 * When using getTitle on every handler, it makes us handle the error every and each time.
 * Wrapping it at makeHandler allows us to overcome this code repetition
 */
// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	match := validPath.FindStringSubmatch(r.URL.Path)
// 	if match == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("Invalid page title")
// 	}
// 	return match[2], nil
// }

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
