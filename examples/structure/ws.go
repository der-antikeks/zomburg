package main

import (
	"io"
	"log"

	"code.google.com/p/go.net/websocket"
)

func wsHandler(ws *websocket.Conn) {
	log.Printf("websocket: %v", ws.RemoteAddr())
	io.Copy(ws, ws)
}
