package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSameKeyIsProcessedOnlyOnceUnderConcurrentRetries(t *testing.T) {
	api := NewAPI()
	server := httptest.NewServer(http.HandlerFunc(api.ingest))
	defer server.Close()
	var group sync.WaitGroup
	for range 20 {
		group.Add(1)
		go func() {
			defer group.Done()
			req, _ := http.NewRequest(http.MethodPost, server.URL, nil)
			req.Header.Set("Idempotency-Key", "same-batch")
			response, err := http.DefaultClient.Do(req)
			if err == nil {
				_ = response.Body.Close()
			}
		}()
	}
	group.Wait()
	if api.processed != 1 {
		t.Fatalf("processed = %d, want 1", api.processed)
	}
}

func TestKeyIsRequired(t *testing.T) {
	api := NewAPI()
	recorder := httptest.NewRecorder()
	api.ingest(recorder, httptest.NewRequest(http.MethodPost, "/ingest", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", recorder.Code)
	}
}
