package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Simple handler that responds to everything
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)
		
		if r.URL.Path == "/health" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"ok"}`)
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"working","message":"proxy active"}`)
		}
	})
	
	// Get port from Railway or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	
	// Start server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
