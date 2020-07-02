package main

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	const project string = "test_dir"

	params := &Parameters{
		ProjectName: project,
	}

	createProject(params)

	config := readConfig(project)
	assert(t, true, config.Build.RemoveBuildDir)
	assert(t, true, config.Build.CopyStaticDirs)
}
