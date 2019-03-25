package main

import (
	"flag"
)

type parameters struct {
	projectName   string
	createProject bool
	buildProject  bool
	cleanProject  bool
}

func readParameters() parameters {
	projectName := flag.String("project", "", "path to the project")
	createProject := flag.Bool("create", false, "set to create new project")
	buildProject := flag.Bool("build", false, "set to build project")
	cleanProject := flag.Bool("clean", false, "set to clean built project")

	flag.Parse()

	return parameters{
		projectName:   *projectName,
		createProject: *createProject,
		buildProject:  *buildProject,
		cleanProject:  *cleanProject,
	}
}
