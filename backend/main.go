package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Xferant VPN Backend", "status": "running"}`)
	})

	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "healthy", "version": "1.0.0", "service": "xferant-vpn"}`)
	})

	http.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"users": []}`)
	})

	http.HandleFunc("/api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "User registered successfully", "token": "jwt-token-placeholder"}`)
	})

	http.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Login successful", "token": "jwt-token-placeholder"}`)
	})

	log.Printf("🚀 Starting Xferant VPN Backend on %s", port)
	log.Printf("📡 Health check: http://localhost%s/api/v1/health", port)
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}