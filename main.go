package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ProxyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	IP      string `json:"ip,omitempty"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := ProxyResponse{
		Status:  "ok",
		Message: "Rota proxy server is running",
	}
	json.NewEncoder(w).Encode(response)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Simple proxy response for now
	response := ProxyResponse{
		Status:  "working",
		Message: "Proxy rotation active",
		IP:      "192.168.1.100", // Will be replaced with real rotation
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Use Railway's PORT environment variable or default to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	
	// Health endpoint
	http.HandleFunc("/health", healthHandler)
	
	// Proxy endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			proxyHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	
	log.Printf("🚀 Rota Proxy Server starting on port %s", port)
	log.Printf("📊 Health: http://localhost:%s/health", port)
	log.Printf("🔄 Proxy: http://localhost:%s/", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
