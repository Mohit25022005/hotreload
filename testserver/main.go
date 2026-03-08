package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Health struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, " HotsgsRel Demo Server from recruiter ")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		resp := Health{
			Status: "ok",
			Time:   time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server running on :9090")

	http.ListenAndServe(":9090", nil)
}