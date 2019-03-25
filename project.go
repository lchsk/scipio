package main

import (
	"os"
	"path/filepath"
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
	err = createDir(filepath.Join(project, "source", "themes"))
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
