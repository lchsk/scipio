package main

import (
	"io/ioutil"
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
	assert(t, checkIfExists(filepath.Join(project, "themes")), true)
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
	filePost, _ := os.Create(filepath.Join(project, "source", "posts", "post_1.md"))
	defer filePost.Close()

	post := `---
title: Article 1
created: 1950-05-15T00:00:00Z
description: Description of Article 1
keywords: scipio, tests, go
tags: tag1, tag2
---

Body of Article 1.

- *Test 1*

- Go

- Scipio
`

	filePost.WriteString(post)

	filePage, _ := os.Create(filepath.Join(project, "source", "pages", "page_1.md"))
	defer filePage.Close()

	page := `---
title: Article 2
created: 1960-05-15T00:00:00Z
description: Description of Article 2
keywords: go, and, stuff
tags: tag1, tag2
---

Body of Article 2.
`

	filePage.WriteString(page)

}

func createTestTheme(project string) {
	filePost, err := os.Create(filepath.Join(project, "themes", "default", "post.html"))
	defer filePost.Close()

	checkError(err)

	template := `<html>
    <head>
        <meta name="description" content="{{description}}">
        <meta name="keywords" content="{{keywords}}">
        <meta name="author" content="">
        <link rel="stylesheet" type="text/css" href="static/styles.css">
    </head>
    <body>
        <div>{{@article-2}}</div>
        <h1>{{title}}</h1>
        {{date}}
        <p>{{body}}</p>
    </body>
</html>`

	filePost.WriteString(template)
}

func TestParseSourceFile(t *testing.T) {
	const project string = "test_dir"

	createProject(project)
	createTestSourceFiles(project)

	data := parseSourceFile(filepath.Join(project, "source", "posts", "post_1.md"), POST)

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

- *Test 1*

- Go

- Scipio`)

	os.RemoveAll(project)
}

func TestBuildProject(t *testing.T) {
	const project string = "test_dir"

	createProject(project)
	createDir(filepath.Join(project, "themes", "default"))
	createTestSourceFiles(project)
	createTestTheme(project)

	buildProject(project)

	html, err := ioutil.ReadFile(filepath.Join(project, "build", "article-1.html"))

	checkError(err)

	expected := `<html>
    <head>
        <meta name="description" content="Description of Article 1">
        <meta name="keywords" content="scipio, tests, go">
        <meta name="author" content="">
        <link rel="stylesheet" type="text/css" href="static/styles.css">
    </head>
    <body>
        <div><a href="article-1.html" title="Description of Article 1">Article 1</a></div>
        <h1>Article 1</h1>
        1950-05-15
        <p><p>Body of Article 1.</p>

<ul>
<li><p><em>Test 1</em></p></li>

<li><p>Go</p></li>

<li><p>Scipio</p></li>
</ul>
</p>
    </body>
</html>`

	assert(t, expected, string(html))

	os.RemoveAll(project)
}
