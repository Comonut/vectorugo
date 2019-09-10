package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/Comonut/vectorugo/store"
)

type controllerConfiguration struct {
	store store.Store
}

func Init(size uint32, name string, persistance bool) {
	// var store, _ = store.NewPersitantStore(size, name+"_index.bin", name+"_vectors.bin")
	var s store.Store
	if persistance {
		s, _ = store.NewPersitantStore(size, name+"_index.bin", name+"_vectors.bin", name+"_search.bin")
	} else {
		s = store.NewSimpleMapStore()
	}
	var config = controllerConfiguration{s}

	http.HandleFunc("/vectors", config.handleVectors)
	http.HandleFunc("/search", config.searchByID)

	logrus.Info("Listening on :8080")
	var err = http.ListenAndServe(":8080", nil)

	if err != nil {
		logrus.Error(err.Error())
	}
}

func (config *controllerConfiguration) handleVectors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key, ok := r.URL.Query()["id"]

		if !ok || len(key[0]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request param 'id' is missing")
			LogRequest(r, http.StatusBadRequest)
			return
		}

		value, err := config.store.Get(key[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Value not present in store: ")
			LogRequest(r, http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error serializing response vector")
			LogRequest(r, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		LogRequest(r, http.StatusOK)

	case "POST":
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error reading request body: ", err)
			LogRequest(r, http.StatusInternalServerError)
			return
		}
		var v map[string][]float64

		err = json.Unmarshal(b, &v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid json format:	", err)
			LogRequest(r, http.StatusBadRequest)
			return
		}
		for k, v := range v {
			err = config.store.Set(k, &store.MemoryVector{ID: k, Array: v})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Error writing to store", err)
				LogRequest(r, http.StatusInternalServerError)
				return
			}
		}
		LogRequest(r, http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid request method - Only GET and POST are supported")
		LogRequest(r, http.StatusMethodNotAllowed)

	}

}

func (config *controllerConfiguration) searchByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key, ok := r.URL.Query()["id"]

		if !ok || len(key[0]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request param 'id' is missing")
			LogRequest(r, http.StatusBadRequest)
			return
		}

		k, ok := r.URL.Query()["k"]
		if !ok || len(k[0]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request param 'k' is missing")
			LogRequest(r, http.StatusBadRequest)
			return
		}

		kN, err := strconv.Atoi(k[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "k has to be an integer ")
			LogRequest(r, http.StatusBadRequest)
			return
		}

		value, err := config.store.Get(key[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Value not present in store: ")
			LogRequest(r, http.StatusBadRequest)
			return
		}

		results, err := config.store.KNN(value, kN)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error getting search results")
			LogRequest(r, http.StatusInternalServerError)
			return
		}
		response := make([]SearchResponseModel, len(*results))
		for i, result := range *results {
			response[i] = SearchResponseModel{ID: result.Target.Name(), Distance: result.Distance}
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error serializing response vector")
			LogRequest(r, http.StatusInternalServerError)
			return
		}
		LogRequest(r, http.StatusOK)

	case "POST":
		k, ok := r.URL.Query()["k"]
		if !ok || len(k[0]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request param 'k' is missing")
			LogRequest(r, http.StatusBadRequest)
			return
		}

		kN, err := strconv.Atoi(k[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "k has to be an integer ")
			LogRequest(r, http.StatusBadRequest)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error reading request body: ", err)
			LogRequest(r, http.StatusBadRequest)
			return
		}
		var v []float64

		err = json.Unmarshal(b, &v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid vector in body", err)
			LogRequest(r, http.StatusBadRequest)
			return
		}

		queryVector := &store.MemoryVector{ID: "", Array: v}
		results, err := config.store.KNN(queryVector, kN)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error getting search results")
			LogRequest(r, http.StatusBadRequest)
			return
		}
		response := make([]SearchResponseModel, len(*results))
		for i, result := range *results {
			response[i] = SearchResponseModel{ID: result.Target.Name(), Distance: result.Distance}
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error serializing response vector")
			LogRequest(r, http.StatusBadRequest)
			return
		}
		LogRequest(r, http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid request method - Only GET is supported")
		LogRequest(r, http.StatusBadRequest)

	}

}
