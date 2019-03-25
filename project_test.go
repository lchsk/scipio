package main

import (
	"testing"
    "path/filepath"
    "os"
)

func TestCheckNilError(t *testing.T) {
    var err error = nil

    checkError(err)
}

func TestCreateProject(t *testing.T) {
    const project string = "test_dir"

    createProject(project)

    assert(t, checkIfExists(project), true)
    assert(t, checkIfExists(filepath.Join(project, "source")), true)
    assert(t, checkIfExists(filepath.Join(project, "source", "posts")), true)
    assert(t, checkIfExists(filepath.Join(project, "source", "pages")), true)
    assert(t, checkIfExists(filepath.Join(project, "source", "themes")), true)
    assert(t, checkIfExists(filepath.Join(project, "source", "index.md")), true)

    assert(t, checkIfExists(filepath.Join(project, "build")), true)

    assert(t, checkIfExists(filepath.Join(project, "scipio.conf")), true)

    os.RemoveAll(project)
}
