package main

import (
	"io/ioutil"
    "log"
    "fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func createProject(project string) {
	var err error

	err = createDir(project)
	checkError(err)

	err = createDir(filepath.Join(project, "source"))
	checkError(err)
	err = createDir(filepath.Join(project, "source", "posts"))
	checkError(err)
	err = createDir(filepath.Join(project, "source", "pages"))
	checkError(err)
	err = createDir(filepath.Join(project, "themes"))
	checkError(err)
	err = createFile(filepath.Join(project, "source", "index.md"))
	checkError(err)

	err = createDir(filepath.Join(project, "build"))
	checkError(err)

	err = createFile(filepath.Join(project, "scipio.conf"))
	checkError(err)
}

func cleanBuild(project string) {
	err := os.RemoveAll(filepath.Join(project, "build"))
	checkError(err)

	err = createDir(filepath.Join(project, "build"))
	checkError(err)
}

const (
	INDEX = iota
	POST
	PAGE
)

type sourceFile struct {
	source      string
	title       string
	description string
	keywords    []string
	tags        []string
	created     time.Time
	body        string
	entryType   int
}

func getValueFromSource(source string, pattern *regexp.Regexp) string {
	matches := pattern.FindStringSubmatch(source)

	if matches != nil && len(matches) == 2 {
		return matches[1]
	}

	return ""
}

func parseSourceFile(project string) sourceFile {
	f, _ := ioutil.ReadFile(filepath.Join(project, "source", "posts", "post_1.md"))
	source := string(f)

	patterns := make(map[string]*regexp.Regexp)
	patterns["title"] = regexp.MustCompile("title: (.+)")
	patterns["description"] = regexp.MustCompile("description: (.+)")
	patterns["keywords"] = regexp.MustCompile("keywords: (.+)")
	patterns["tags"] = regexp.MustCompile("tags: (.+)")
	patterns["created"] = regexp.MustCompile("created: (.+)")
	patterns["body"] = regexp.MustCompile("(?s)---.*---(.*)")

	title := getValueFromSource(source, patterns["title"])
	description := getValueFromSource(source, patterns["description"])
	body := strings.TrimSpace(getValueFromSource(source, patterns["body"]))

	parsedKeywords := strings.Split(getValueFromSource(source, patterns["keywords"]), ",")
	keywords := make([]string, len(parsedKeywords))

	for i, keyword := range parsedKeywords {
		keywords[i] = strings.TrimSpace(keyword)
	}

	parsedTags := strings.Split(getValueFromSource(source, patterns["tags"]), ",")
	tags := make([]string, len(parsedTags))

	for i, tag := range parsedTags {
		tags[i] = strings.TrimSpace(tag)
	}

	t, err := time.Parse(time.RFC3339, getValueFromSource(source, patterns["created"]))

	if err != nil {
		t = time.Unix(0, 0)
	}

	data := sourceFile{
		source:      source,
		title:       title,
		description: description,
		keywords:    keywords,
		tags:        tags,
		created:     t,
		body:        body,
		entryType:   POST,
	}

	return data
}

func buildProject(project string) {
	files, err := ioutil.ReadDir(filepath.Join(project, "source", "posts"))

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(files[0].Name())

	parseSourceFile("test_dir")
}
