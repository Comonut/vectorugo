package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Comonut/vectorugo/store"
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

		var v store.Vector
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, &v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid json format:	", err)
			return
		}
		fmt.Fprint(w, v.Sum())

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "ERROR")

	}
}

func searchByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Recieved search with request\n", *r)
}
