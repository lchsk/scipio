package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	Url    string
	OutputExtension string `toml:"output_extension"`
	LinksBeginWithSlash bool `toml:"links_begin_with_slash"`
	Rss    rssConfig
	Posts  postsConfig
	Static staticConfig
	Build  buildConfig
}

type staticConfig struct {
	Copy []map[string]string
}

type buildConfig struct {
	RemoveBuildDir bool `toml:"remove_build_dir"`
	CopyStaticDirs bool `toml:"copy_static_dirs"`
}

type rssConfig struct {
	GenerateRss bool `toml:"generate_rss"`
	Title       string
	Description string
	AuthorName  string `toml:"author_name"`
	AuthorEmail string `toml:"author_email"`
}

type postsConfig struct {
	Slug        string   `toml:slug`
	Title       string   `toml:title`
	Description string   `toml:description`
	Keywords    []string `toml:keywords`
}

func readConfig(project string) *config {
	configPath := filepath.Join(project, "scipio.toml")
	f, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatal("Config file scipio.toml does not exist ", err)
	}

	conf := &config{}
	// Defaults
	conf.Build.RemoveBuildDir = true
	conf.Build.CopyStaticDirs = true
	conf.OutputExtension = ".html"
	conf.LinksBeginWithSlash = false

	if _, err := toml.Decode(string(f), conf); err != nil {
		log.Fatal("Failed to read the config file! ", err)
	}

	return conf

}
