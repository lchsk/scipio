package main

import (
	"fmt"
	"os"
)

func main() {
	params := readParameters()

	if params.projectName == "" {
		fmt.Println("--project must not be empty")
		os.Exit(1)
	}

	if params.createProject {
		createProject(params.projectName)
		os.Exit(0)
	}

	if params.cleanBuild {
		cleanBuild(params.projectName)
		os.Exit(0)
	}

	if params.buildProject {
		buildProject(params.projectName)
		os.Exit(0)
	}
}
