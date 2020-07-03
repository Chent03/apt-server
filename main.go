package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	http.ListenAndServe(":3000", r)
}
