package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Proxy list - add your real proxies here
var proxies = []string{
	// Add your real HTTP proxies here:
	// "http://username:password@proxy1.example.com:8080",
	// "http://username:password@proxy2.example.com:8080", 
	// "http://username:password@proxy3.example.com:8080",
	
	// Premium residential proxies (best for Telegram):
	// "http://user:pass@residential.provider.com:8080",
	
	// Datacenter proxies (good for testing):
	// "http://user:pass@datacenter.provider.com:8080",
	
	// Working free proxies (test these first):
	"http://45.8.179.247:8080",
	"http://45.130.40.133:8080",
	"http://185.199.229.156:3128",
	
	// If none work, use direct connection (no proxy)
	"DIRECT",
}

var currentProxyIndex = 0
var proxyMutex sync.Mutex

func getNextProxy() string {
	proxyMutex.Lock()
	defer proxyMutex.Unlock()
	
	proxy := proxies[currentProxyIndex]
	currentProxyIndex = (currentProxyIndex + 1) % len(proxies)
	return proxy
}

func handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Proxy request: %s %s from %s", r.Method, r.URL.Path, r.Host)
	
	// Get next proxy from rotation
	proxyURL := getNextProxy()
	log.Printf("Using proxy: %s", proxyURL)
	
	var transport *http.Transport
	
	// Handle DIRECT connection (no proxy)
	if proxyURL == "DIRECT" {
		log.Printf("Using direct connection (no proxy)")
		transport = &http.Transport{}
	} else {
		// Parse proxy URL
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			log.Printf("Invalid proxy URL: %s", err)
			http.Error(w, "Invalid proxy configuration", http.StatusInternalServerError)
			return
		}
		
		// Create proxy transport
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}
	
	// Create new request with the same method and headers
	targetURL := r.URL
	if !strings.HasPrefix(targetURL.String(), "http") {
		targetURL.Scheme = "http"
		targetURL.Host = r.Host
	}
	
	req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	
	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	
	// Make request through proxy
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Add timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Proxy request failed: %s", err)
		http.Error(w, fmt.Sprintf("Proxy request failed: %s", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	
	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	
	// Copy response status and body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	
	log.Printf("Proxy response: %d %s", resp.StatusCode, resp.Status)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"proxy rotation active","proxies":%d}`, len(proxies))
}

func main() {
	// Get port from Railway or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Health endpoint
	http.HandleFunc("/health", healthHandler)
	
	// Proxy handler for all other requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			healthHandler(w, r)
			return
		}
		handleProxyRequest(w, r)
	})
	
	log.Printf(" Rota Proxy Server starting on port %s", port)
	log.Printf(" Health check: http://localhost:%s/health", port)
	log.Printf(" Proxy rotation with %d proxies", len(proxies))
	
	// Start server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
