package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (this *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.once.Do(func() {
		this.templ = template.Must(template.ParseFiles(filepath.Join("templates", this.fileName)))
	})
	this.templ.Execute(w, nil)
}

func main() {
	fmt.Println("vim-go")
	t := templateHandler{fileName: "chat.html"}
	http.Handle("/", &t)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
