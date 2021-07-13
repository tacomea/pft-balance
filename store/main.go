package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
}
