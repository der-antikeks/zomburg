package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("listening %v", ":8080")
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(""))))
}
