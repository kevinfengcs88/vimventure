package main

// CHANGE THIS LATER AND IMPORT INTO THE ACTUAL MAIN PROGRAM

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "this is the home response")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	http.ListenAndServe(":6969", nil)
}
