package main

import "fmt"

type Signals struct {
	P95Milliseconds int
	CPUThrottled    bool
	ConsumerLag     int
}

func diagnose(s Signals) string {
	if s.ConsumerLag > 1000 && (s.P95Milliseconds > 500 || s.CPUThrottled) {
		return "action: bottleneck confirmed; protect traffic and inspect CPU limits/downstream latency"
	}
	if s.ConsumerLag > 1000 {
		return "action: backlog detected; investigate consumer throughput before scaling"
	}
	return "action: no correlated backend bottleneck"
}

func main() {
	scenarios := []Signals{
		{P95Milliseconds: 80, ConsumerLag: 40},
		{P95Milliseconds: 780, CPUThrottled: true, ConsumerLag: 8200},
		{P95Milliseconds: 100, ConsumerLag: 5000},
	}
	for _, scenario := range scenarios {
		fmt.Printf("%+v -> %s\n", scenario, diagnose(scenario))
	}
}
