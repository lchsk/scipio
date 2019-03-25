package main

import (
	"fmt"
)

func main() {
	params := ReadParameters()

	fmt.Println(params.ProjectName)
}
