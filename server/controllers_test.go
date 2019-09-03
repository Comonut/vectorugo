package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Comonut/vectorugo/store"
)

func TestMain(t *testing.T) {
	os.Remove("test_index.bin")
	os.Remove("test_vectors.bin")
	os.Remove("test_search.bin")
	go Init(uint32(4), "test", false)
	time.Sleep(2 * time.Second)

	testSetGet(t)
	testSearchController(t)
}

func testSetGet(t *testing.T) {

	resp, err := http.Post("http://localhost:8080/vectors", "application/json", bytes.NewBuffer([]byte("{\"v1\" asdasd: [0, 0.0, 1, 3.14]}")))
	if err != nil || resp.StatusCode != 400 {
		t.Errorf("Expected bad request for malformed json")
	}
	resp, err = http.Get("http://localhost:8080/vectors")
	if err != nil || resp.StatusCode != 400 {
		t.Errorf("Expected 400 for vector request without specified ID")
	}
	resp, err = http.Get("http://localhost:8080/vectors?id=v1")
	if err != nil || resp.StatusCode != 404 {
		t.Errorf("Vector was not set but didn't get 404")
	}
	resp, err = http.Post("http://localhost:8080/vectors", "application/json", bytes.NewBuffer([]byte("{\"v1\" : [0, 0.0, 1, 3.14]}")))
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("Error setting values %d", resp.StatusCode)
	}
	resp, err = http.Get("http://localhost:8080/vectors?id=v1")
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("Error getting values %d", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error("could not read response")
	}
	var v store.MemoryVector
	if json.Unmarshal(b, &v) != nil {
		t.Error("could not parse response from get")
	}

	expected := store.MemoryVector{ID: "v1", Array: []float64{0, 0, 1, 3.14}}

	if v.ID != expected.ID || len(*expected.Values()) != len(*v.Values()) {
		t.Error("Received wrong vector lol")
	}

	for i := range *v.Values() {
		if (*v.Values())[i] != (*expected.Values())[i] {
			t.Errorf("diffent values at position %d - expected %f , but got %f", i, (*expected.Values())[i], (*v.Values())[i])
		}
	}

}

func testSearchController(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/vectors", "application/json", bytes.NewBuffer([]byte("{\"v1\" : [0, 0.0, 1, 3.14]}")))
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("Error setting values %d", resp.StatusCode)
	}

	resp, err = http.Get("http://localhost:8080/search?id=v1")
	if err != nil || resp.StatusCode != 400 {
		t.Errorf("No K value defined - expected 400")
	}

	resp, err = http.Get("http://localhost:8080/search?k=5")
	if err != nil || resp.StatusCode != 400 {
		t.Errorf("No target defined - expected 400")
	}
	resp, err = http.Get("http://localhost:8080/search?id=v1&k=0.5")
	if err != nil || resp.StatusCode != 400 {
		t.Errorf("Accepted a non-integer K value - expected 400")
	}

	resp, err = http.Get("http://localhost:8080/search?id=v2&k=1")
	if err != nil || resp.StatusCode != 404 {
		t.Errorf("Tried KNN on non-existing vector - expected 404")
	}

	resp, err = http.Get("http://localhost:8080/search?id=v1&k=1")
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("Error getting values %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error("could not read response")
	}
	var v []SearchResponseModel
	if json.Unmarshal(b, &v) != nil {
		t.Error("could not parse response from get")
	}

	if v[0].ID != "v1" || v[0].Distance != 0 {
		t.Error("wrong result")
	}

}
