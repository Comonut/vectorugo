package server

import (
	"net/http"
	"testing"
	"time"
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
