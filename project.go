package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/avelino/slugify"
	"github.com/gorilla/feeds"
	"gopkg.in/russross/blackfriday.v2"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func createProject(project string, params *Parameters) {
	var err error

	err = createDir(project)
	checkError(err)

	err = createDir(filepath.Join(project, "source"))
	checkError(err)
	err = createDir(filepath.Join(project, "source", "posts"))
	checkError(err)
	err = createDir(filepath.Join(project, "source", "pages"))
	checkError(err)
	err = createDir(filepath.Join(project, "source", "data"))
	checkError(err)
	err = createDir(filepath.Join(project, "themes"))
	checkError(err)
	err = createDir(filepath.Join(project, "themes", "default"))
	checkError(err)
	err = createDir(filepath.Join(project, "themes", "default", "static"))
	checkError(err)
	err = createFile(filepath.Join(project, "themes", "default", "index.html"))
	checkError(err)
	err = createFile(filepath.Join(project, "source", "index.md"))
	checkError(err)

	err = createFile(filepath.Join(project, "scipio.toml"))
	checkError(err)
}

func getBuildDir(project string, params *Parameters) string {
	if params.BuildDir == "" {
		return filepath.Join(project, "build")
	} else {
		return params.BuildDir
	}
}

func cleanBuild(project string, params *Parameters) {
	err := os.RemoveAll(getBuildDir(project, params))
	checkError(err)

	err = createDir(getBuildDir(project, params))
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

func getArticles(files []os.FileInfo, postsDir string, entryType int) []sourceFile {
	var articles []sourceFile

	for _, postFile := range files {
		name := postFile.Name()
		ext := name[len(name)-3:]
		if ext == ".md" {
			data := parseSourceFile(filepath.Join(postsDir, name), entryType)
			articles = append(articles, data)
		}
	}

	return articles
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
	articles []sourceFile, conf *config, templates map[string]string, params *Parameters) {
	themeFilePath := filepath.Join(project, "themes", theme, templateFile)
	themeHtml, err := ioutil.ReadFile(themeFilePath)

	checkError(err)

	output := string(themeHtml)

	buildDir := getBuildDir(project, params)

	outputFilePath := filepath.Join(buildDir, data.slug+".html")
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

	for templateName, templateContent := range templates {
		if templateFile == templateName {
			continue
		}
		output = strings.Replace(output, "{{#include "+templateName+"}}", templateContent, -1)
	}

	extraPath := regexp.MustCompile("{{#include (.+)}}")

	matches := extraPath.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		path := filepath.Join(project, "source", match[1])
		f, err := ioutil.ReadFile(path)

		if err != nil {
			continue
		}

		output = strings.Replace(output, match[0], string(f), -1)
	}

	for _, article := range articles {
		output = strings.Replace(output, "{{@"+article.slug+"}}", createLink(article), -1)
	}

	if templateFile == "index.html" || templateFile == "posts.html" {
		output = addPostsLinksHtml(output, project, "default", templateFile, data, articles, conf, params)
	}

	os.Truncate(outputFilePath, 0)
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

func generateRss(project string, articles []sourceFile, conf *config, params *Parameters) {
	buildDir := getBuildDir(project, params)
	rssPath := filepath.Join(buildDir, "posts.xml")
	outputFile, err := os.OpenFile(rssPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer outputFile.Close()

	checkError(err)

	var items []*feeds.Item

	for _, article := range articles {
		link := &feeds.Link{Href: fmt.Sprintf("%s/%s.html", conf.Url, article.slug)}
		items = append(items, &feeds.Item{Title: article.title, Link: link, Description: article.description, Created: article.created})
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       conf.Rss.Title,
		Link:        &feeds.Link{Href: conf.Url},
		Description: conf.Rss.Description,
		Author:      &feeds.Author{Name: conf.Rss.AuthorName, Email: conf.Rss.AuthorEmail},
		Created:     now,
		Items:       items,
	}

	rss, err := feed.ToRss()

	if err != nil {
		log.Fatal(err)
	}

	os.Truncate(rssPath, 0)
	outputFile.WriteString(rss)
}

func addPostsLinksHtml(output string, project string, theme string, templateFile string, data sourceFile,
	articles []sourceFile, conf *config, params *Parameters) string {

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

		if conf.Rss.GenerateRss {
			generateRss(project, sortedPosts, conf, params)
		}
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

func buildProject(project string, conf *config, params *Parameters) {
	cleanBuild(project, params)
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

	templatesDir := filepath.Join(project, "themes", "default")
	templateFiles, err := ioutil.ReadDir(templatesDir)

	if err != nil {
		log.Fatal(err)
	}

	templates := make(map[string]string)

	for _, templateFile := range templateFiles {
		f, _ := ioutil.ReadFile(filepath.Join(templatesDir, templateFile.Name()))
		templates[templateFile.Name()] = string(f)
	}

	postsArticles := getArticles(postsFiles, postsDir, POST)
	pagesArticles := getArticles(pagesFiles, pagesDir, PAGE)

	articles := append(postsArticles, pagesArticles...)

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

		generateArticleHtml(project, "default", template, article, articles, conf, templates, params)
	}

	// TODO: Move this to the config file
	postArticle := sourceFile{
		slug:        conf.Posts.Slug,
		title:       conf.Posts.Title,
		description: conf.Posts.Description,
		keywords:    conf.Posts.Keywords,
	}

	// Don't generate if Posts config is not defined
	if conf.Posts.Slug != "" {
		generateArticleHtml(project, "default", "posts.html", postArticle, articles, conf, templates, params)
	}

	copyStaticDirectories(project, conf, params)
}

func copyStaticDirectories(project string, conf *config, params *Parameters) {
	staticFrom := filepath.Join(project, "themes", "default", "static")
	buildDir := getBuildDir(project, params)
	dataFrom := filepath.Join(project, "source", "data")

	copyDirectory(staticFrom, buildDir)
	copyDirectory(dataFrom, buildDir)

	for _, dirData := range conf.Static.Copy {
		from := filepath.Join(project, "source", dirData["from"])
		to := filepath.Join(buildDir, dirData["to"])
		copyDirectory(from, to)
	}
}
