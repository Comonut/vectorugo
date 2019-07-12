package server

import (
	"fmt"
	"net/http"
)

func Init() {
	http.HandleFunc("/vectors", handleVectors)
	http.HandleFunc("/search", searchByID)

	var err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Print(err.Error())
	}
}

func handleVectors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprint(w, "GET")
	case "POST":
		fmt.Fprint(w, "POST")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "ERROR")

	}
}

func searchByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Recieved search with request\n", *r)
}
