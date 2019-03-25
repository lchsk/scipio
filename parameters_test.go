package main

import (
	"testing"
)

func TestReadParameters(t *testing.T) {
	params := ReadParameters()

	assert(t, params.ProjectName, "")
	assert(t, params.CreateProject, false)
	assert(t, params.BuildProject, false)
	assert(t, params.CleanProject, false)
}
