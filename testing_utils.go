package main

import (
	"testing"
)

func assertWithMessage(t *testing.T, variable interface{}, value interface{}, message string) {
	if variable != value {
		t.Errorf("Expected '%v', got '%v' %s", variable, value, message)
	}
}

func assert(t *testing.T, variable interface{}, value interface{}) {
	assertWithMessage(t, variable, value, "")
}
