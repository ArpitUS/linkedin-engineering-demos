package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1ContractRemainsStable(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/v1/customers/c-1", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	var actual map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	if actual["name"] != "Ada Lovelace" {
		t.Fatalf("name = %v", actual["name"])
	}
	if _, hasTier := actual["tier"]; hasTier {
		t.Fatal("v1 unexpectedly changed")
	}
}

func TestV2IsAdditive(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/v2/customers/c-1", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	var actual CustomerV2
	if err := json.Unmarshal(recorder.Body.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	if actual.Tier != "pro" || actual.Name != "Ada Lovelace" {
		t.Fatalf("payload = %+v", actual)
	}
}
