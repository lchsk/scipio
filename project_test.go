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

	params := &Parameters{
		ProjectName: project,
	}

	createProject(params)

	assert(t, checkIfExists(project), true)
	assert(t, checkIfExists(filepath.Join(project, "source")), true)
	assert(t, checkIfExists(filepath.Join(project, "source", "posts")), true)
	assert(t, checkIfExists(filepath.Join(project, "source", "pages")), true)
	assert(t, checkIfExists(filepath.Join(project, "source", "data")), true)
	assert(t, checkIfExists(filepath.Join(project, "themes")), true)
	assert(t, checkIfExists(filepath.Join(project, "themes", "default", "static")), true)
	assert(t, checkIfExists(filepath.Join(project, "source", "index.md")), true)

	assert(t, checkIfExists(filepath.Join(project, "build")), false)

	assert(t, checkIfExists(filepath.Join(project, "scipio.toml")), true)

	os.RemoveAll(project)
}

func getDefaultConfig() *config {
	build := buildConfig{RemoveBuildDir: true}
	return &config{Build: build, OutputExtension: ".html", LinksBeginWithSlash: false}
}

func TestCleanBuild(t *testing.T) {
	const project string = "test_dir"

	params := &Parameters{
		ProjectName: project,
	}

	createProject(params)
	buildProject(getDefaultConfig(), params)
	createFile(filepath.Join(project, "build", "test_file"))

	assert(t, checkIfExists(filepath.Join(project, "build", "test_file")), true)

	cleanBuild(params)

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

	fileIndex, _ := os.Create(filepath.Join(project, "source", "index.md"))
	defer fileIndex.Close()

	index := `---
title: Index Page
created: 1970-05-15T00:00:00Z
description: Index Page Description
keywords: go, and, stuff
tags: tag1, tag2
---

Home Page
`

	fileIndex.WriteString(index)
}

func createTestTheme(project string) {
	fileTopMenuTemplate, err := os.Create(filepath.Join(project, "themes", "default", "top_menu.html"))
	defer fileTopMenuTemplate.Close()
	checkError(err)

	filePost, err := os.Create(filepath.Join(project, "themes", "default", "post.html"))
	defer filePost.Close()
	checkError(err)

	filePostList, err := os.Create(filepath.Join(project, "themes", "default", "posts.html"))
	defer filePostList.Close()
	checkError(err)

	fileTopMenuTemplate.WriteString("top menu")

	postTemplate := `<html>
    <head>
        <meta name="description" content="{{description}}">
        <meta name="keywords" content="{{keywords}}">
        <meta name="author" content="">
        <link rel="stylesheet" type="text/css" href="static/styles.css">
    </head>
    <body>
        {{#include top_menu.html}}
        <div>{{@article-2}}</div>
        <h1>{{title}}</h1>
        {{date}}
        <p>{{body}}</p>
    </body>
</html>`

	filePost.WriteString(postTemplate)

	filePage, errPage := os.Create(filepath.Join(project, "themes", "default", "page.html"))
	defer filePage.Close()

	checkError(errPage)

	pageTemplate := `<html>
    <head>
        <meta name="description" content="{{description}}">
        <meta name="keywords" content="{{keywords}}">
        <meta name="author" content="">
        <link rel="stylesheet" type="text/css" href="static/styles.css">
    </head>
    <body>
        {{#include top_menu.html}}
        <div>{{@article-2}}</div>
        <h1>{{title}}</h1>
        {{date}}
        <p>{{body}}</p>
    </body>
</html>`

	filePage.WriteString(pageTemplate)

	fileIndex, errIndex := os.Create(filepath.Join(project, "themes", "default", "index.html"))
	defer fileIndex.Close()

	checkError(errIndex)

	indexTemplate := `<html>
    <head>
        <meta name="description" content="{{description}}">
        <meta name="keywords" content="{{keywords}}">
        <meta name="author" content="">
        <link rel="stylesheet" type="text/css" href="static/styles.css">
    </head>
    <body>
        {{#include top_menu.html}}
        <div>{{@article-2}}</div>
        <h1>{{title}}</h1>
        {{date}}
        <p>{{body}}</p>
{{posts-begin}}
        <li>{{post_link}}</li>
        <span class="date">{{post_date}}</span>
{{posts-end}}
    </body>
</html>`

	fileIndex.WriteString(indexTemplate)
}

func TestParseSourceFile(t *testing.T) {
	const project string = "test_dir"

	params := &Parameters{
		ProjectName: project,
	}

	createProject(params)
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

	params := &Parameters{
		ProjectName: project,
	}

	createProject(params)
	createDir(filepath.Join(project, "themes", "default"))
	createTestSourceFiles(project)
	createTestTheme(project)

	buildProject(getDefaultConfig(), params)

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
        top menu
        <div><a href="article-2.html" title="Description of Article 2">Article 2</a></div>
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

	indexHtml, indexErr := ioutil.ReadFile(filepath.Join(project, "build", "index.html"))

	checkError(indexErr)

	expectedIndex := `<html>
    <head>
        <meta name="description" content="Index Page Description">
        <meta name="keywords" content="go, and, stuff">
        <meta name="author" content="">
        <link rel="stylesheet" type="text/css" href="static/styles.css">
    </head>
    <body>
        top menu
        <div><a href="article-2.html" title="Description of Article 2">Article 2</a></div>
        <h1>Index Page</h1>
        1970-05-15
        <p><p>Home Page</p>
</p>

        <li><a href="welcome-to-your-new-website.html" title="Description of the first post">Welcome to your new website</a></li>
        <span class="date">2021-01-27</span>

        <li><a href="article-1.html" title="Description of Article 1">Article 1</a></li>
        <span class="date">1950-05-15</span>

    </body>
</html>`

	assert(t, expectedIndex, string(indexHtml))

	os.RemoveAll(project)
}
