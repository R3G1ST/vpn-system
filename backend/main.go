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
		fmt.Fprintf(w, "Xferant VPN Backend")
	})
	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","version":"1.0.0"}`)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})
	log.Printf("Server starting on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}