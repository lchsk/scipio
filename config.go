package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"path/filepath"
)

type config struct {
	Url string
	Rss rssConfig
}

type rssConfig struct {
	GenerateRss bool `toml:"generate_rss"`
	Title       string
	Description string
	AuthorName  string `toml:"author_name"`
	AuthorEmail string `toml:"author_email"`
}

func readConfig(project string) *config {
	configPath := filepath.Join(project, "scipio.toml")
	f, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatal("Config file scipio.toml does not exist ", err)
	}

	conf := &config{}

	if _, err := toml.Decode(string(f), conf); err != nil {
		log.Fatal("Failed to read the config file! ", err)
	}

	return conf

}
