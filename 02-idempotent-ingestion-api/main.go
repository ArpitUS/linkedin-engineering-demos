package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
)

type receipt struct {
	Status string `json:"status"`
	ID     int    `json:"id"`
}

type API struct {
	mu        sync.Mutex
	byKey     map[string]receipt
	nextID    int
	processed int
}

func NewAPI() *API { return &API{byKey: make(map[string]receipt)} }

func (a *API) ingest(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		http.Error(w, "Idempotency-Key is required", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if prior, exists := a.byKey[key]; exists {
		writeJSON(w, http.StatusOK, prior)
		return
	}
	a.nextID++
	a.processed++ // Represents one irreversible database or billing write.
	result := receipt{Status: "accepted", ID: a.nextID}
	a.byKey[key] = result
	writeJSON(w, http.StatusCreated, result)
}

func writeJSON(w http.ResponseWriter, status int, value receipt) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func main() {
	api := NewAPI()
	server := httptest.NewServer(http.HandlerFunc(api.ingest))
	defer server.Close()
	for _, attempt := range []int{1, 2} {
		req, _ := http.NewRequest(http.MethodPost, server.URL, nil)
		req.Header.Set("Idempotency-Key", "batch-2026-09-04-a")
		response, _ := http.DefaultClient.Do(req)
		var body receipt
		_ = json.NewDecoder(response.Body).Decode(&body)
		_ = response.Body.Close()
		fmt.Printf("attempt %d: HTTP %d, receipt=%+v\n", attempt, response.StatusCode, body)
	}
	fmt.Printf("irreversible writes: %d\n", api.processed)
}
