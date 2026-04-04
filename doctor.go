package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func checkCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func checkFile(base, name string) bool {
	path := filepath.Join(base, name)
	_, err := os.Stat(path)
	return err == nil
}

func runDoctor(base string) ([]string, bool) {
	var result []string
	ok := true

	files := []string{
		"main.go",
		"icon.png",
		"go.mod",
		"go.sum",
	}

	for _, f := range files {
		if checkFile(base, f) {
			result = append(result, "✅ "+f)
		} else {
			result = append(result, "❌ "+f+" missing")
			ok = false
		}

		if ok {
			result = append(result, "✅ All checks passed")
		}

	}

	return result, ok
}
