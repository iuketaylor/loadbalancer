package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.String("port", "8081", "Server port")
var serverId = flag.String("id", "1", "Server Id")

func handleHome(w http.ResponseWriter, req *http.Request) {
	log.Printf("Replied with a hello message")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello From Backend Server %s\n", *serverId)
}

func main() {
	flag.Parse()

	log.Printf("Backend server %s listening on port %s...", *serverId, *port)
	http.HandleFunc("/", handleHome)
	http.ListenAndServe(":"+*port, nil)
}
