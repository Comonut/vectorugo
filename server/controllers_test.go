package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Comonut/vectorugo/store"
)

func TestRequests(t *testing.T) {
	go Init()
	time.Sleep(2 * time.Second)

	var resp, _ = http.Get("http://localhost:8080/get")
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 for /get bur received %d", resp.StatusCode)
	}

	resp, _ = http.Get("http://localhost:8080/set")
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 for /set bur received %d", resp.StatusCode)
	}

	resp, _ = http.Get("http://localhost:8080/search")
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 for /search bur received %d", resp.StatusCode)
	}

	resp, _ = http.Get("http://localhost:8080/nonexistant")
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404 for /search bur received %d", resp.StatusCode)
	}
}

func TestSetting(t *testing.T) {

	// go Init()
	time.Sleep(2 * time.Second)

	values := make(map[string]*store.Vector)

	for q := 0; q < 10; q++ {
		values[fmt.Sprintf("Vector%d", q)] = store.Random(fmt.Sprintf("Vector%d", q), 1024)
	}

	byts, _ := json.Marshal(values)
	resp, err := http.Post("http://localhost:8080/vectors", "application/json", bytes.NewBuffer(byts))
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("Error setting values %d", resp.StatusCode)
	}

	fmt.Print("done")
}
