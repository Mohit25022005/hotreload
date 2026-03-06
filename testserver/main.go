package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Teri maa ka bhosda")
	})

	fmt.Println("Server running on :9090")

	http.ListenAndServe(":9090", nil)
}
