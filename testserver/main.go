package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/users", usersHandler)
	mux.HandleFunc("/slow", slowHandler)

	fmt.Println("Server running on :" + port)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "🔥 HotReload Live Reload Works!\n")
	fmt.Fprintf(w, "Time: %s\n", time.Now().Format(time.RFC3339))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	resp := HealthResponse{
		Status: "ok",
		Time:   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	users := []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
		{"id": "3", "name": "Charlie"},
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {

	time.Sleep(2 * time.Second)

	fmt.Fprintln(w, "Slow endpoint finished")
}

func loggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		log.Printf("➡️  %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("⬅️  completed in %v", time.Since(start))
	})
}