package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	servers []string
	current int
	mutex   sync.Mutex
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		current: 0,
	}
}

func (lb *LoadBalancer) getNextServer() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	server := lb.servers[lb.current]

	// increment selected server then wrap back round to 0 when reached total servers (round robin)
	lb.current = (lb.current + 1) % len(lb.servers)

	log.Printf("Selected server: %s", server)
	return server
}

var lb *LoadBalancer

func handleHome(w http.ResponseWriter, req *http.Request) {
	serverUrl := lb.getNextServer()

	forwardRequest(w, req, serverUrl)
}

func forwardRequest(w http.ResponseWriter, req *http.Request, serverUrl string) {
	backendReq, err := http.NewRequest(req.Method, serverUrl, req.Body)
	if err != nil {
		http.Error(w, "Failed to create backend request", http.StatusInternalServerError)
		return
	}

	for key, values := range req.Header {
		for _, value := range values {
			backendReq.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(backendReq)
	if err != nil {
		http.Error(w, "Backend request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Printf("Response from server %s: %s %s", serverUrl, resp.Proto, resp.Status)

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading backend response: %v", err)
		return
	}

	log.Printf("Response body: %s", string(body))
	w.Write(body)
}

func main() {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	lb = NewLoadBalancer(servers)

	log.Println("Load balancer starting on port 8080...")
	log.Printf("List of servers: %v", servers)

	http.HandleFunc("/", handleHome)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
