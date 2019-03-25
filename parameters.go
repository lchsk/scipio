package main

import (
	"flag"
)

type Parameters struct {
	ProjectName   string
	CreateProject bool
	BuildProject  bool
	CleanProject  bool
}

func ReadParameters() Parameters {
	projectName := flag.String("project", "", "path to the project")
	createProject := flag.Bool("create", false, "set to create new project")
	buildProject := flag.Bool("build", false, "set to build project")
	cleanProject := flag.Bool("clean", false, "set to clean built project")

	flag.Parse()

	return Parameters{
		ProjectName:   *projectName,
		CreateProject: *createProject,
		BuildProject:  *buildProject,
		CleanProject:  *cleanProject,
	}
}
