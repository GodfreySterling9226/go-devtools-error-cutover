package main

import "fmt"

type BuildEvent struct {
	ID          string            `json:"-"`
	Title       string            `json:"title"`
	Message     string            `json:"message"`
	Level       string            `json:"level"`
	Fingerprint []string          `json:"fingerprint"`
	Exception   string            `json:"exception"`
	Context     map[string]string `json:"context"`
}

type Release struct {
	Version string
	Commit  string
}

func captureBuildFailure(client *InfraiClient, release Release, stage string, cause error) error {
	event := BuildEvent{ID: release.Version + "-" + stage, Title: "build stage failed", Message: cause.Error(), Level: "error", Fingerprint: []string{"build", stage}, Exception: fmt.Sprintf("%T: %v", cause, cause), Context: map[string]string{"release": release.Version, "commit": release.Commit, "stage": stage}}
	return client.Capture(event)
}

func NextAction(buildOK, captureOK bool) string {
	if buildOK {
		return "promote release"
	}
	if captureOK {
		return "rollback release"
	}
	return "halt and page on-call"
}
