package main

import "fmt"

type Event struct {
	ID      string
	Service string
	Value   int
}

type NaiveStore struct{ events []Event }

func (s *NaiveStore) Ingest(event Event) { s.events = append(s.events, event) }

type DeduplicatingStore struct{ events map[string]Event }

func NewDeduplicatingStore() *DeduplicatingStore {
	return &DeduplicatingStore{events: make(map[string]Event)}
}

// Ingest returns whether this is the first successful delivery of the event.
func (s *DeduplicatingStore) Ingest(event Event) bool {
	if _, exists := s.events[event.ID]; exists {
		return false
	}
	s.events[event.ID] = event
	return true
}

func main() {
	event := Event{ID: "evt-9f2", Service: "checkout", Value: 1}

	naive := &NaiveStore{}
	naive.Ingest(event) // Server writes event, but the acknowledgement is lost.
	naive.Ingest(event) // Client retries because it never observed success.
	fmt.Printf("naive pipeline: %d events, total=%d (incorrect)\n", len(naive.events), sum(naive.events))

	safe := NewDeduplicatingStore()
	fmt.Printf("first delivery accepted: %t\n", safe.Ingest(event))
	fmt.Printf("retry accepted: %t\n", safe.Ingest(event))
	fmt.Printf("idempotent pipeline: %d event, total=%d (correct)\n", len(safe.events), event.Value)
}

func sum(events []Event) int {
	total := 0
	for _, event := range events {
		total += event.Value
	}
	return total
}
