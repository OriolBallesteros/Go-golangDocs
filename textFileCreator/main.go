package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	title, body := provideFlags()

	page1 := &Page{
		*title,
		[]byte(*body),
	}
	page1.saveInFile()

	page2, _ := loadFile(page1.Title)
	fmt.Println(string(page2.Body))
}

func provideFlags() (*string, *string) {
	title := flag.String("wt", "title", "activate this flag to write directly on the console the title of your txt file")
	body := flag.String("wb", "This is a sample Page body.", "activate this flag to write directly on the console the body of your txt file")
	flag.Parse()

	return title, body
}

func (p *Page) saveInFile() error {
	fileName := p.Title + ".txt"

	return ioutil.WriteFile(fileName, p.Body, 0600)
}

func loadFile(title string) (*Page, error) {
	fileName := title + ".txt"
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}
