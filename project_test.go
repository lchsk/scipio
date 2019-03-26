package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestCleanBuild(t *testing.T) {
	const project string = "test_dir"

	createProject(project)
	createFile(filepath.Join(project, "build", "test_file"))

	assert(t, checkIfExists(filepath.Join(project, "build", "test_file")), true)

	cleanBuild(project)

	assert(t, checkIfExists(filepath.Join(project, "build", "test_file")), false)

	os.RemoveAll(project)
}

func createTestSourceFiles(project string) {
	f, _ := os.Create(filepath.Join(project, "source", "posts", "post_1.md"))
	defer f.Close()

	post := `---
title: Article 1
created: 1950-05-15T00:00:00Z
description: Description of Article 1
keywords: scipio, tests, go
tags: tag1, tag2
---

Body of Article 1.

- Test 1

- Go

- Scipio
`

	f.WriteString(post)
}

func TestParseSourceFile(t *testing.T) {
	const project string = "test_dir"

	createProject(project)
	createTestSourceFiles(project)

	data := parseSourceFile(project)

	assert(t, data.title, "Article 1")
	assert(t, data.description, "Description of Article 1")

	assert(t, data.keywords[0], "scipio")
	assert(t, data.keywords[1], "tests")
	assert(t, data.keywords[2], "go")
	assert(t, len(data.keywords), 3)

	assert(t, data.tags[0], "tag1")
	assert(t, data.tags[1], "tag2")
	assert(t, len(data.tags), 2)

	assert(t, data.created, time.Date(1950, time.May, 15, 0, 0, 0, 0, time.UTC))
	assert(t, data.entryType, POST)

	assert(t, data.body, `Body of Article 1.

- Test 1

- Go

- Scipio`)

	os.RemoveAll(project)
}

func TestBuildProject(t *testing.T) {
	const project string = "test_dir"

	createProject(project)
	createTestSourceFiles(project)

	buildProject(project)

	os.RemoveAll(project)
}
