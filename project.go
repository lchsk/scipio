package main

import (
	"fmt"
	"github.com/avelino/slugify"
	"github.com/gorilla/feeds"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	source       string
	title        string
	description  string
	keywords     []string
	tags         []string
	redirects    []string
	created      time.Time
	body         string
	entryType    int
	slug         string
	overrideSlug string
}

func getValueFromSource(source string, pattern *regexp.Regexp) string {
	matches := pattern.FindStringSubmatch(source)

	if matches != nil && len(matches) == 2 {
		return matches[1]
	}

	return ""
}

func parseSourceFile(path string, entryType int) sourceFile {
	f, _ := ioutil.ReadFile(path)
	source := string(f)

	patterns := make(map[string]*regexp.Regexp)
	patterns["title"] = regexp.MustCompile("title: (.+)")
	patterns["overrideSlug"] = regexp.MustCompile("slug: (.+)")
	patterns["description"] = regexp.MustCompile("description: (.+)")
	patterns["keywords"] = regexp.MustCompile("keywords: (.+)")
	patterns["tags"] = regexp.MustCompile("tags: (.+)")
	patterns["redirectFrom"] = regexp.MustCompile("redirect_from: (.+)")
	patterns["created"] = regexp.MustCompile("created: (.+)")
	patterns["body"] = regexp.MustCompile("(?s)---.*---(.*)")

	title := getValueFromSource(source, patterns["title"])
	overrideSlug := getValueFromSource(source, patterns["overrideSlug"])
	description := getValueFromSource(source, patterns["description"])
	body := strings.TrimSpace(getValueFromSource(source, patterns["body"]))

	// TODO: Factor this out
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

	parsedRedirects := strings.Split(getValueFromSource(source, patterns["redirectFrom"]), ",")
	var redirects []string

	for _, redirect := range parsedRedirects {
		trimmedRedirect := strings.TrimSpace(redirect)

		if trimmedRedirect != "" {
			redirects = append(redirects, trimmedRedirect)
		}
	}

	t, err := time.Parse(time.RFC3339, getValueFromSource(source, patterns["created"]))

	if err != nil {
		t = time.Unix(0, 0)
	}

	slug := ""

	if entryType == INDEX {
		slug = "index"
	} else {
		if overrideSlug == "" {
			slug = slugify.Slugify(title)
		} else {
			slug = overrideSlug
		}
	}

	data := sourceFile{
		source:      source,
		title:       title,
		description: description,
		keywords:    keywords,
		tags:        tags,
		redirects:   redirects,
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

	if len(data.redirects) != 0 {
		generateRedirectHtml(project, data)
	}

	output = strings.Replace(output, "{{title}}", data.title, -1)
	output = strings.Replace(output, "{{description}}", data.description, -1)
	output = strings.Replace(output, "{{body}}", string(blackfriday.Run([]byte(data.body))), -1)
	output = strings.Replace(output, "{{keywords}}", strings.Join(data.keywords, ", "), -1)
	output = strings.Replace(output, "{{date}}", data.created.Format("2006-01-02"), -1)

	for _, article := range articles {
		output = strings.Replace(output, "{{@"+article.slug+"}}", createLink(article), -1)
	}

	if templateFile == "index.html" {
		output = addIndexHtml(output, project, "default", templateFile, data, articles)
	}

	outputFile.WriteString(output)
}

func filterArticles(articles []sourceFile, entryType int) []sourceFile {
	var filtered []sourceFile

	for _, article := range articles {
		if article.entryType == entryType {
			filtered = append(filtered, article)
		}
	}

	return filtered
}

func sortArticles(articles []sourceFile) []sourceFile {
	// Sort in descending order by date
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].created.After(articles[j].created)
	})

	return articles
}

func generateRss(project string, articles []sourceFile) {
	outputFile, err := os.OpenFile(filepath.Join(project, "build", "posts.xml"), os.O_CREATE|os.O_WRONLY, 0644)
	defer outputFile.Close()

	checkError(err)

	var items []*feeds.Item

	for _, article := range articles {
		link := &feeds.Link{Href: fmt.Sprintf("https://lchsk.com/%s.html", article.slug)}
		items = append(items, &feeds.Item{Title: article.title, Link: link, Description: article.description, Created: article.created})
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "lchsk",
		Link:        &feeds.Link{Href: "http://lchsk.com"},
		Description: "lchsk",
		Author:      &feeds.Author{Name: "lchsk", Email: "lchsk"},
		Created:     now,
		Items:       items,
	}

	rss, err := feed.ToRss()

	if err != nil {
		log.Fatal(err)
	}

	outputFile.WriteString(rss)
}

func addIndexHtml(output string, project string, theme string, templateFile string, data sourceFile,
	articles []sourceFile) string {

	sortedPosts := filterArticles(sortArticles(articles), POST)

	postsLoop := regexp.MustCompile("(?s){{posts-begin}}(.*){{posts-end}}")
	posts := postsLoop.FindStringSubmatch(output)

	if len(posts) == 2 {
		postsContent := ""

		for _, post := range sortedPosts {
			singlePostContent := strings.Replace(posts[1], "{{post_link}}", createLink(post), -1)
			singlePostContent = strings.Replace(singlePostContent, "{{post_date}}", post.created.Format("2006-01-02"), -1)
			postsContent += singlePostContent
		}

		output = strings.Replace(output, posts[1], postsContent, -1)
		output = strings.Replace(output, "{{posts-begin}}", "", -1)
		output = strings.Replace(output, "{{posts-end}}", "", -1)

		generateRss(project, sortedPosts)
	}

	return output
}

func generateRedirectHtml(project string, data sourceFile) {
	for _, redirect := range data.redirects {
		outputFilePath := filepath.Join(project, "build", redirect+".html")
		outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644)

		defer outputFile.Close()

		checkError(err)

		os.Truncate(outputFilePath, 0)

		newUrl := data.slug + ".html"

		redirectHtml := fmt.Sprintf(`<html>
<head>
<meta http-equiv="refresh" content="0; url=%s">
<body>
<p><a href="%s">Redirect</a></p>
</body>
</html>
`, newUrl, newUrl)
		outputFile.WriteString(redirectHtml)
	}
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
	}

	copyStaticDirectories(project)
}

func copyStaticDirectories(project string) {
	staticFrom := filepath.Join(project, "themes", "default", "static")
	staticTo := filepath.Join(project, "build")

	dataFrom := filepath.Join(project, "source", "data")
	dataTo := filepath.Join(project, "build")

	copyDirectory(staticFrom, staticTo)
	copyDirectory(dataFrom, dataTo)
}
