package server

import (
	"fmt"
	"net/http"
)

func Init() {
	http.HandleFunc("/get", getById)
	http.HandleFunc("/set", setById)
	http.HandleFunc("/search", searchById)
	http.ListenAndServe(":8080", nil)
}

func getById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Recieved get with request\n%s!", r)
}

func setById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Recieved set with request\n%s", r)
}

func searchById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Recieved search with request\n%s", r)
}
