package main

import (
	"os"
	"os/exec"
)

func checkIfExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func createDir(dir string) error {
	if checkIfExists(dir) {
		return nil
	}

	return os.Mkdir(dir, os.ModePerm)
}

func createFile(path string) error {
	if checkIfExists(path) {
		return nil
	}

	f, err := os.Create(path)
	defer f.Close()

	return err
}

func copyDirectory(from, to string) {
	cmd := exec.Command("cp", "-r", from, to)
	_, err := cmd.Output()

	if err != nil {
		panic(err)
	}
}
