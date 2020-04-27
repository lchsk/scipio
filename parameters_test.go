package main

import (
	"testing"
)

func TestReadParameters(t *testing.T) {
	params := readParameters()

	assert(t, params.ProjectName, "")
	assert(t, params.CreateProject, false)
	assert(t, params.BuildProject, false)
	assert(t, params.CleanBuild, false)
	assert(t, params.Serve, false)
}
