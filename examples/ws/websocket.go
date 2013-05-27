package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go.net/websocket"
)

type Message struct {
	Message string
	Time    time.Time
}

type Server struct {
	connections map[*User]bool
	broadcast   chan Message
	register    chan *User
	unregister  chan *User
}

func (s *Server) Run() {
	for {
		select {
		case u := <-s.register:
			log.Printf("register: %v", u)

			s.connections[u] = true
			go func() { s.broadcast <- Message{"User entered the room", time.Now()} }()
		case u := <-s.unregister:
			log.Printf("unregister: %v", u)

			delete(s.connections, u)
			close(u.send)

			go func() { s.broadcast <- Message{"User left the room", time.Now()} }()
		case m := <-s.broadcast:
			log.Printf("broadcast: %v", m)

			for u := range s.connections {
				select {
				case u.send <- m:
				default:
					delete(s.connections, u)
					close(u.send)
					go u.ws.Close()
				}
			}
		}
	}
}

var chat = Server{
	connections: make(map[*User]bool),
	broadcast:   make(chan Message),
	register:    make(chan *User),
	unregister:  make(chan *User),
}

type User struct {
	ws   *websocket.Conn
	send chan Message
}

func (u *User) read() {
	log.Printf("%v waiting for read", u)

	for {
		var m Message
		if websocket.JSON.Receive(u.ws, &m) != nil {
			break
		}
		log.Printf("%v received: %v", u, m)
		chat.broadcast <- m
	}

	u.ws.Close()
	log.Printf("%v closed read", u)
}

func (u *User) write() {
	log.Printf("%v waiting for write", u)

	for m := range u.send {
		log.Printf("%v sending: %v", u, m)

		if websocket.JSON.Send(u.ws, m) != nil {
			break
		}
	}

	u.ws.Close()
	log.Printf("%v closed write", u)
}

var webListen, wsAddr string

//var tmpl *template.Template
//no template caching for testing purposes only

func init() {
	flag.StringVar(&webListen, "listen", ":8080", "address to listen for HTTP/WebSockets on")
	flag.StringVar(&wsAddr, "ws", "localhost:8080", "websocket host[:port], as seen by JavaScript")

}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("html: %v", r.URL)

	u := "A"

	tmpl := template.Must(template.ParseFiles("template.html"))
	err := tmpl.Execute(w, map[string]interface{}{
		"WSAddr": wsAddr,
		"User":   u,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func streamHandler(ws *websocket.Conn) {
	log.Printf("websocket: %v", ws.RemoteAddr())

	u := &User{
		send: make(chan Message, 256),
		ws:   ws,
	}

	log.Printf("register new: %v", u)

	chat.register <- u
	defer func() {
		chat.unregister <- u
	}()

	log.Printf("start loop for: %v", u)
	go u.write()
	u.read()
}

func main() {
	flag.Parse()

	go chat.Run()

	http.HandleFunc("/", htmlHandler)
	http.Handle("/stream", websocket.Handler(streamHandler))

	log.Printf("starting server at: %v", webListen)
	log.Printf("waiting for websocket at: %v", wsAddr)
	log.Fatal(http.ListenAndServe(webListen, nil))
}
