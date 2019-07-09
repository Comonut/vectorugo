package server

import (
	"fmt"
	"net/http"
)

func Init() {
	http.HandleFunc("/get", getByID)
	http.HandleFunc("/set", setByID)
	http.HandleFunc("/search", searchByID)
	http.ListenAndServe(":8080", nil)
}

func getByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Recieved get with request\n!", *r)
}

func setByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Recieved set with request\n", *r)
}

func searchByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Recieved search with request\n", *r)
}
