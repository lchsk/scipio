package main

import (
	"os"
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
