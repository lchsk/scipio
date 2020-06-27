package main

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func setupReloading(conf *config, params *Parameters) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("File %s was modified, rebuilding...", event.Name)
					buildProject(conf, params)
					log.Printf("Project was rebuilt")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	project := params.ProjectName

	dirs := []string{}
	dirs = append(dirs, filepath.Join(project, "scipio.toml"))
	dirs = append(dirs, filepath.Join(project, "source"))
	dirs = append(dirs, filepath.Join(project, "source", "pages"))
	dirs = append(dirs, filepath.Join(project, "source", "posts"))
	dirs = append(dirs, filepath.Join(project, "themes", "default"))
	dirs = append(dirs, filepath.Join(project, "themes", "default", "static"))

	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
}
