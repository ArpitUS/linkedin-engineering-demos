package main

import "testing"

func TestRetryCreatesDuplicateWithoutEventIdentity(t *testing.T) {
	store := &NaiveStore{}
	event := Event{ID: "stable-id", Value: 1}
	store.Ingest(event)
	store.Ingest(event)
	if got := sum(store.events); got != 2 {
		t.Fatalf("naive total = %d, want 2", got)
	}
}

func TestStableEventIDMakesRetrySafe(t *testing.T) {
	store := NewDeduplicatingStore()
	event := Event{ID: "stable-id", Value: 1}
	if !store.Ingest(event) || store.Ingest(event) {
		t.Fatal("expected first delivery only to be accepted")
	}
	if got := len(store.events); got != 1 {
		t.Fatalf("stored events = %d, want 1", got)
	}
}
