package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	http.ListenAndServe(":"+port, r)
}
