package main

import (
	"fmt"
	"net/http"
	"os"
)

const ScipioVersion = "0.1"

func main() {
	params := readParameters()

	if params.Version {
		fmt.Printf("Scipio v%s\n", ScipioVersion)
		os.Exit(0)
	}

	if params.ProjectName == "" {
		fmt.Println("--project must not be empty")
		os.Exit(1)
	}

	if params.CreateProject {
		createProject(params)
		os.Exit(0)
	}

	if params.CleanBuild {
		cleanBuild(params)
		os.Exit(0)
	}

	conf := readConfig(params.ProjectName)

	if params.BuildProject {
		buildProject(conf, params)
		os.Exit(0)
	}

	if params.Serve {
		buildProject(conf, params)
		http.Handle("/", http.FileServer(http.Dir(getBuildDir(params))))
		fmt.Println("Serving content on http://localhost:4000")
		if err := http.ListenAndServe(":4000", nil); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}
