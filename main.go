package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	
	log.Println("🚀 Rota Proxy Server starting on :8000")
	log.Println("📊 Health: http://localhost:8000/health")
	log.Println("🔄 Proxy: http://localhost:8000/")
	
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
# In your rota-proxy-server repository:
