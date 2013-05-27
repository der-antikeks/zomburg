package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

var webListen, wsAddr string
var tmpl *template.Template

func init() {
	flag.StringVar(&webListen, "listen", ":8080", "address to listen for HTTP on")
	flag.StringVar(&wsAddr, "ws", "localhost:8080", "websocket host[:port], as seen by JavaScript")

	tmpl = template.Must(template.ParseFiles("html/template.html"))
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("incoming: %v > %v", r.RemoteAddr, r.URL)

	err := tmpl.Execute(w, map[string]interface{}{
		"WSAddr": wsAddr,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", htmlHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/geo/", http.StripPrefix("/geo/", http.FileServer(http.Dir("./geo"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./img/favicon.ico")
	})

	log.Printf("starting server at: %v", webListen)
	log.Fatal(http.ListenAndServe(webListen, nil))
}
