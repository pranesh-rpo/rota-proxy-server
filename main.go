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
		log.Printf("Request received: %s %s from %s", r.Method, r.URL.Path, r.Host)
		
		// Always return 200 OK for health checks
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		if r.URL.Path == "/health" {
			fmt.Fprintf(w, `{"status":"ok","message":"healthy"}`)
		} else {
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
	log.Printf("🔍 Accepting requests from healthcheck.railway.app")
	
	// Start server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
