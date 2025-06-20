package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleHome(w http.ResponseWriter, req *http.Request) {
	log.Printf("Received request from %s\n", req.RemoteAddr)
	log.Printf("%s %s %s", req.Method, req.URL.Path, req.Proto)
	log.Printf("Host: %s\n", req.Host)
	log.Printf("User-Agent: %s\n", req.UserAgent())
	log.Printf("Accept %s\n", req.Header.Get("Accept"))

	fmt.Fprintf(w, "Hello\n")
}
func main() {
	http.HandleFunc("/", handleHome)

	log.Println("Listening on port :8080")
	http.ListenAndServe(":8080", nil)
}
