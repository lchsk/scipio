package main

import (
	"flag"
)

type Parameters struct {
	ProjectName   string
	CreateProject bool
	BuildProject  bool
	CleanBuild    bool
	Serve         bool
}

func readParameters() *Parameters {
	// TODO: Change projectName variable name
	projectName := flag.String("project", "", "path to the project")
	createProject := flag.Bool("create", false, "set to create new project")
	buildProject := flag.Bool("build", false, "build project and quit")
	cleanBuild := flag.Bool("clean", false, "set to clean built project")
	serve := flag.Bool("serve", false, "build and run server")

	flag.Parse()

	return &Parameters{
		ProjectName:   *projectName,
		CreateProject: *createProject,
		BuildProject:  *buildProject,
		CleanBuild:    *cleanBuild,
		Serve:         *serve,
	}
}
