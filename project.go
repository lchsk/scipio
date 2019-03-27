package main

import (
	"fmt"
	"github.com/avelino/slugify"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
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
	slug        string
}

func getValueFromSource(source string, pattern *regexp.Regexp) string {
	matches := pattern.FindStringSubmatch(source)

	if matches != nil && len(matches) == 2 {
		return matches[1]
	}

	return ""
}

func parseSourceFile(path string, entryType int) sourceFile {
	// f, _ := ioutil.ReadFile(filepath.Join(project, "source", "posts", "post_1.md"))
	f, _ := ioutil.ReadFile(path)
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

	slug := ""

	if entryType == INDEX {
		slug = "index"
	} else {
		slug = slugify.Slugify(title)
	}

	data := sourceFile{
		source:      source,
		title:       title,
		description: description,
		keywords:    keywords,
		tags:        tags,
		created:     t,
		body:        body,
		entryType:   entryType,
		slug:        slug,
	}

	return data
}

func generateArticleHtml(project string, theme string, templateFile string, data sourceFile,
	articles []sourceFile) {
	themeFilePath := filepath.Join(project, "themes", theme, templateFile)
	themeHtml, err := ioutil.ReadFile(themeFilePath)

	checkError(err)

	output := string(themeHtml)

	outputFilePath := filepath.Join(project, "build", data.slug+".html")
	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	defer outputFile.Close()

	checkError(err)

	output = strings.Replace(output, "{{title}}", data.title, -1)
	output = strings.Replace(output, "{{description}}", data.description, -1)
	output = strings.Replace(output, "{{body}}", string(blackfriday.Run([]byte(data.body))), -1)
	output = strings.Replace(output, "{{keywords}}", strings.Join(data.keywords, ", "), -1)
	output = strings.Replace(output, "{{date}}", data.created.Format("2006-01-02"), -1)

	for _, article := range articles {
		output = strings.Replace(output, "{{@"+article.slug+"}}", createLink(article), -1)
	}

	outputFile.WriteString(output)
}

func generateIndexHtml(project string, theme string, templateFile string, data sourceFile,
	articles []sourceFile) {

}

func createLink(data sourceFile) string {
	return fmt.Sprintf("<a href=\"%s\" title=\"%s\">%s</a>", data.slug+".html", data.description, data.title)
}

func buildProject(project string) {
	postsDir := filepath.Join(project, "source", "posts")
	postsFiles, err := ioutil.ReadDir(postsDir)

	if err != nil {
		log.Fatal(err)
	}

	pagesDir := filepath.Join(project, "source", "pages")
	pagesFiles, err := ioutil.ReadDir(pagesDir)

	if err != nil {
		log.Fatal(err)
	}

	var articles []sourceFile

	for _, postFile := range postsFiles {
		data := parseSourceFile(filepath.Join(postsDir, postFile.Name()), POST)
		articles = append(articles, data)
	}

	for _, pageFile := range pagesFiles {
		data := parseSourceFile(filepath.Join(pagesDir, pageFile.Name()), PAGE)
		articles = append(articles, data)
	}

	indexData := parseSourceFile(filepath.Join(project, "source", "index.md"), INDEX)
	articles = append(articles, indexData)

	for _, article := range articles {
		template := ""

		if article.entryType == POST {
			template = "post.html"
		} else if article.entryType == PAGE {
			template = "page.html"
		} else {
			template = "index.html"
		}

		generateArticleHtml(project, "default", template, article, articles)

		if template == "index.html" {
			generateIndexHtml(project, "default", template, article, articles)
		}
	}
}
