package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
)

func requestHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

func main() {
	flag.Parse()
	user := template.Must(template.ParseFiles("index.html"))
	admin := template.Must(template.ParseFiles("admin.html"))
	host := template.Must(template.ParseFiles("host.html"))
	h := newHub()
	router := http.NewServeMux()
	router.Handle("/admin", requestHandler(admin))
	router.Handle("/host", requestHandler(host))
	router.Handle("/static/", scriptHandler())
	router.Handle("/", requestHandler(user))
	router.Handle("/customer_ws", customerWSHandler{h: h})
	router.Handle("/admin_ws", adminWSHandler{h: h})
	log.Fatal(http.ListenAndServe(":8081", router))
}

func scriptHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/app.js")
	})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/app.js")
}
