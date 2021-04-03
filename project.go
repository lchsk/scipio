package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/avelino/slugify"
	"github.com/gorilla/feeds"
	"github.com/russross/blackfriday/v2"
)

var codeBlockRegex = regexp.MustCompile("(?s)```([a-z]+)(\\n)(.+?)```")

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func writeToNewFile(path string, contents string) {
	outputFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	defer outputFile.Close()

	outputFile.WriteString(contents)
}

func createProject(params *Parameters) {
	project := params.ProjectName
	var err error

	err = createDir(project)
	checkError(err)

	err = createDir(filepath.Join(project, "source"))
	checkError(err)
	err = createDir(filepath.Join(project, "source", "posts"))
	checkError(err)

	firstPost := fmt.Sprintf("%s\n\n%s\n", firstPostPattern, "```python\nprint('hello')\n```")

	writeToNewFile(filepath.Join(project, "source", "posts", "welcome.md"), firstPost)
	checkError(err)
	err = createDir(filepath.Join(project, "source", "pages"))
	checkError(err)

	writeToNewFile(filepath.Join(project, "source", "pages", "privacy-policy.md"), privacyPolicy)
	writeToNewFile(filepath.Join(project, "source", "index.md"), indexPage)

	checkError(err)
	err = createDir(filepath.Join(project, "source", "data"))
	checkError(err)
	err = createDir(filepath.Join(project, "themes"))
	checkError(err)
	err = createDir(filepath.Join(project, "themes", "default"))
	checkError(err)
	err = createDir(filepath.Join(project, "themes", "default", "static"))
	checkError(err)

	writeToNewFile(filepath.Join(project, "themes", "default", "index.html"), indexTheme)
	writeToNewFile(filepath.Join(project, "themes", "default", "post.html"), postTheme)
	writeToNewFile(filepath.Join(project, "themes", "default", "page.html"), pageTheme)
	writeToNewFile(filepath.Join(project, "themes", "default", "top.html"), topTheme)
	writeToNewFile(filepath.Join(project, "themes", "default", "footer.html"), footerTheme)
	writeToNewFile(filepath.Join(project, "themes", "default", "header.html"), headerTheme)
	writeToNewFile(filepath.Join(project, "themes", "default", "posts.html"), postsTheme)
	writeToNewFile(filepath.Join(project, "themes", "default", "app.scss"), themeStyleApp)
	writeToNewFile(filepath.Join(project, "themes", "default", "bootstrap.scss"), themeStyleBootstrap)
	writeToNewFile(filepath.Join(project, "themes", "default", "bundle.scss"), themeStyleBundle)
	
	// Scipio.toml config
	writeToNewFile(filepath.Join(project, "scipio.toml"), configValue)

	// package.json
	writeToNewFile(filepath.Join(project, "package.json"), packageJsonValue)
}

func getBuildDir(params *Parameters) string {
	if params.BuildDir == "" {
		return filepath.Join(params.ProjectName, "build")
	} else {
		return params.BuildDir
	}
}

func cleanBuild(params *Parameters) {
	err := os.RemoveAll(getBuildDir(params))
	checkError(err)

	err = createDir(getBuildDir(params))
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
	templateFile string
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
	patterns["templateFile"] = regexp.MustCompile("template: (.+)")
	patterns["description"] = regexp.MustCompile("description: (.+)")
	patterns["keywords"] = regexp.MustCompile("keywords: (.+)")
	patterns["tags"] = regexp.MustCompile("tags: (.+)")
	patterns["redirectFrom"] = regexp.MustCompile("redirect_from: (.+)")
	patterns["created"] = regexp.MustCompile("created: (.+)")
	patterns["body"] = regexp.MustCompile("(?s)---.*---(.*)")

	title := getValueFromSource(source, patterns["title"])
	overrideSlug := getValueFromSource(source, patterns["overrideSlug"])
	templateFile := getValueFromSource(source, patterns["templateFile"])
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

	if templateFile == "" {
		if entryType == POST {
			templateFile = "post.html"
		} else if entryType == PAGE {
			templateFile = "page.html"
		} else {
			templateFile = "index.html"
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
		templateFile: templateFile,
	}

	return data
}

func generateArticleHtml(theme string, templateFile string, data sourceFile,
	articles []sourceFile, conf *config, templates map[string]string, params *Parameters) {
	project := params.ProjectName
	themeFilePath := filepath.Join(project, "themes", theme, templateFile)
	themeHtml, err := ioutil.ReadFile(themeFilePath)

	checkError(err)

	output := string(themeHtml)

	buildDir := getBuildDir(params)

	var outputExtension string

	if templateFile == "index.html" {
		outputExtension = ".html"
	} else {
		outputExtension = conf.OutputExtension
	}

	outputFilePath := filepath.Join(buildDir, data.slug + outputExtension)
	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	defer outputFile.Close()

	checkError(err)

	if len(data.redirects) != 0 {
		generateRedirectHtml(params, data, conf)
	}

	// Syntax highlighting
	// TODO: Move somewhere else

	//codeBlock := regexp.MustCompile("(?s)```([a-z]+)(\\n)(.+?)```")
	body := data.body
	codeMatches := codeBlockRegex.FindAllStringSubmatch(body, -1)

	for _, match := range codeMatches {
		tmpFilePath := filepath.Join("scipio_pygments_input")
		os.Remove(tmpFilePath)
		o1, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Printf("Could not create a temp file: %s\n", err)
			continue
		}

		if len(match) != 4 {
			fmt.Printf("Couldn't find code block for syntax highlighting, got %d elements instead of 4", len(match))
			continue
		}

		o1.WriteString(match[3])

		language := match[1]

		// TODO: Add optional -O linenos option to add line numbers

		out, err := exec.Command("pygmentize", "-f", "html", "-l", language, "scipio_pygments_input").Output()

		if err != nil {
			fmt.Printf("Couldn't find pygmentize, syntax highlighting will not work: %s\n", err)
			continue
		}

		body = strings.Replace(body, match[0], string(out), -1)
	}

	output = strings.Replace(output, "{{title}}", data.title, -1)
	output = strings.Replace(output, "{{description}}", data.description, -1)
	output = strings.Replace(output, "{{body}}", string(blackfriday.Run([]byte(body))), -1)
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
		output = strings.Replace(output, "{{@"+article.slug+"}}", createLink(article, conf), -1)
	}

	// TODO: template file names shouldn't be hardcoded
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

func generateRss(articles []sourceFile, conf *config, params *Parameters) {
	buildDir := getBuildDir(params)
	rssPath := filepath.Join(buildDir, "posts.xml")
	outputFile, err := os.OpenFile(rssPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer outputFile.Close()

	checkError(err)

	var items []*feeds.Item
	outputExtension := conf.OutputExtension

	for _, article := range articles {
		link := &feeds.Link{Href: fmt.Sprintf("%s/%s%s", conf.Url, article.slug, outputExtension)}
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
			singlePostContent := strings.Replace(posts[1], "{{post_link}}", createLink(post, conf), -1)
			singlePostContent = strings.Replace(singlePostContent, "{{post_date}}", post.created.Format("2006-01-02"), -1)
			singlePostContent = strings.Replace(singlePostContent, "{{post_description}}", post.description, -1)
			postsContent += singlePostContent
		}

		output = strings.Replace(output, posts[1], postsContent, -1)
		output = strings.Replace(output, "{{posts-begin}}", "", -1)
		output = strings.Replace(output, "{{posts-end}}", "", -1)

		if conf.Rss.GenerateRss {
			generateRss(sortedPosts, conf, params)
		}
	}

	return output
}

func generateRedirectHtml(params *Parameters, data sourceFile, conf *config) {
	outputExtension := conf.OutputExtension

	for _, redirect := range data.redirects {
		outputFilePath := filepath.Join(getBuildDir(params), redirect + outputExtension)
		outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644)

		defer outputFile.Close()

		checkError(err)

		os.Truncate(outputFilePath, 0)

		newUrl := data.slug + outputExtension

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

func createLink(data sourceFile, conf *config) string {
	outputExtension := conf.OutputExtension
	linksBeginWithSlash := conf.LinksBeginWithSlash

	var linkStart string
	if linksBeginWithSlash {
		linkStart = "/"
	} else {
		linkStart = ""
	}

	return fmt.Sprintf("<a href=\"%s%s\" title=\"%s\">%s</a>", linkStart, data.slug + outputExtension, data.description, data.title)
}

func buildProject(conf *config, params *Parameters) {
	project := params.ProjectName

	if conf.Build.RemoveBuildDir {
		cleanBuild(params)
	}

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
		generateArticleHtml("default", article.templateFile, article, articles, conf, templates, params)
	}

	// TODO: Move this to the config file
	postArticle := sourceFile{
		slug:        conf.Posts.Slug,
		title:       conf.Posts.Title,
		description: conf.Posts.Description,
		keywords:    conf.Posts.Keywords,
	}

	// Don't generate if Posts config is not defined
	// TODO: templateFile shouldn't be hardecoded
	if conf.Posts.Slug != "" {
		generateArticleHtml("default", "posts.html", postArticle, articles, conf, templates, params)
	}

	if conf.Build.CopyStaticDirs {
		copyStaticDirectories(conf, params)
	}
}

func copyStaticDirectories(conf *config, params *Parameters) {
	project := params.ProjectName
	staticFrom := filepath.Join(project, "themes", "default", "static")
	buildDir := getBuildDir(params)
	dataFrom := filepath.Join(project, "source", "data")

	copyDirectory(staticFrom, buildDir)
	copyDirectory(dataFrom, buildDir)

	for _, dirData := range conf.Static.Copy {
		from := filepath.Join(project, "source", dirData["from"])
		to := filepath.Join(buildDir, dirData["to"])
		copyDirectory(from, to)
	}
}
