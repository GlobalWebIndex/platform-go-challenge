package main

import (
	"io"
	"log"
	"testing"
)

// silenceLogger is a test helper to discard logger output from the test log.
func silenceLogger(t *testing.T) {
	t.Helper()
	prev := log.Writer()
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(prev) })
}
