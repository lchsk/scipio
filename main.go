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
}
