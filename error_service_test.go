package main

import "testing"

func TestNextAction(t *testing.T) {
	tests := []struct {
		name           string
		build, capture bool
		want           string
	}{
		{"green", true, true, "promote release"},
		{"captured failure", false, true, "rollback release"},
		{"unreported failure", false, false, "halt and page on-call"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextAction(tt.build, tt.capture); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
