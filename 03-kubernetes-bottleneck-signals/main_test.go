package main

import (
	"strings"
	"testing"
)

func TestCorrelatedSignalsConfirmBottleneck(t *testing.T) {
	got := diagnose(Signals{P95Milliseconds: 900, CPUThrottled: true, ConsumerLag: 3000})
	if !strings.Contains(got, "confirmed") {
		t.Fatalf("diagnosis = %q", got)
	}
}

func TestLagAloneIsNotEnoughToNameRootCause(t *testing.T) {
	got := diagnose(Signals{P95Milliseconds: 80, ConsumerLag: 3000})
	if !strings.Contains(got, "backlog") {
		t.Fatalf("diagnosis = %q", got)
	}
}
