package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
)

func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

func main() {
	flag.Parse()
	tpl := template.Must(template.ParseFiles("index.html"))
	admin := template.Must(template.ParseFiles("admin.html"))
	host := template.Must(template.ParseFiles("host.html"))
	h := newHub()
	router := http.NewServeMux()
	router.Handle("/", homeHandler(tpl))
	router.Handle("/ws", wsHandler{h: h})
	router.Handle("/admin", homeHandler(admin))
	router.Handle("/admin_ws", adminWSHandler{h: h})
	router.Handle("/host", homeHandler(host))
	log.Printf("serving on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
