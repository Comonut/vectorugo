package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Comonut/vectorugo/store"
)

type controllerConfiguration struct {
	store store.Store
}

func Init(size uint32, name string) {
	// var store, _ = store.NewPersitantStore(size, name+"_index.bin", name+"_vectors.bin")
	s := store.NewSimpleMapStore()
	var config = controllerConfiguration{s}

	http.HandleFunc("/vectors", config.handleVectors)
	http.HandleFunc("/search", config.searchByID)

	var err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Print(err.Error())
	}
}

func (config *controllerConfiguration) handleVectors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key, ok := r.URL.Query()["id"]

		if !ok || len(key[0]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request param 'id' is missing")
			return
		}

		value, err := config.store.Get(key[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Value not present in store: ")
			return
		}
		err = json.NewEncoder(w).Encode(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error serializing response vector")
			return
		}
		w.Header().Set("Content-Type", "application/json")

	case "POST":
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error reading request body: ", err)
			return
		}
		var v map[string][]float64

		err = json.Unmarshal(b, &v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid json format:	", err)
			return
		}
		for k, v := range v {
			err = config.store.Set(k, &store.Vector{ID: k, Values: v})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Error writing to store", err)
				return
			}
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid request method - Only GET and POST are supported")

	}
}

func (config *controllerConfiguration) searchByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key, ok := r.URL.Query()["id"]

		if !ok || len(key[0]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request param 'id' is missing")
			return
		}

		k, ok := r.URL.Query()["k"]
		if !ok || len(k[0]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request param 'k' is missing")
			return
		}

		kN, err := strconv.Atoi(k[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "k has to be an integer ")
			return
		}

		value, err := config.store.Get(key[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Value not present in store: ")
			return
		}

		results, err := config.store.KNN(value, kN)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error getting search results")
			return
		}
		response := make([]SearchResponseModel, len(*results))
		for i, result := range *results {
			response[i] = SearchResponseModel{ID: result.Target.ID, Distance: result.Distance}
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error serializing response vector")
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid request method - Only GET is supported")

	}
}
