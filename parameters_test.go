package main

import (
	"testing"
)

func TestReadParameters(t *testing.T) {
	params := readParameters()

	assert(t, params.projectName, "")
	assert(t, params.createProject, false)
	assert(t, params.buildProject, false)
	assert(t, params.cleanProject, false)
}
