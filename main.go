package main

import (
	"fmt"
	"log"
)

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	release := Release{Version: "2026.09.0", Commit: "abc123"}
	buildErr := fmt.Errorf("compiler: package analytics/loader failed")
	if err := captureBuildFailure(client, release, "compile", buildErr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(NextAction(false, true))
}
