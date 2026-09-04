package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type CustomerV1 struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CustomerV2 adds fields without changing the meaning of V1 fields.
type CustomerV2 struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Tier        string `json:"tier"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/v1/customers/c-1":
		_ = json.NewEncoder(w).Encode(CustomerV1{ID: "c-1", Name: "Ada Lovelace"})
	case "/v2/customers/c-1":
		_ = json.NewEncoder(w).Encode(CustomerV2{ID: "c-1", Name: "Ada Lovelace", DisplayName: "Ada", Tier: "pro"})
	default:
		http.NotFound(w, r)
	}
}

func main() {
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	for _, path := range []string{"/v1/customers/c-1", "/v2/customers/c-1"} {
		response, _ := http.Get(server.URL + path)
		var payload map[string]any
		_ = json.NewDecoder(response.Body).Decode(&payload)
		_ = response.Body.Close()
		fmt.Printf("%s -> %v\n", path, payload)
	}
}
