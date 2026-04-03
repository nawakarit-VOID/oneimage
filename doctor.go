package main

import (
	"os"
	"os/exec"
)

func checkCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func checkFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runDoctor() ([]string, bool) {
	var result []string
	allPassed := true

	if checkCmd("go") {
		result = append(result, "✅ Go installed")
	} else {
		result = append(result, "❌ Go missing")
		allPassed = false
	}

	if checkCmd("fyne") {
		result = append(result, "✅ fyne installed")
	} else {
		result = append(result, "❌ fyne missing")
		allPassed = false
	}

	if checkFile("icon.png") {
		result = append(result, "✅ icon.png found")
	} else {
		result = append(result, "❌ icon.png missing")
		allPassed = false
	}

	if checkFile("main.go") {
		result = append(result, "✅ main.go found")
	} else {
		result = append(result, "❌ main.go missing")
		allPassed = false
	}

	if checkFile("appimagetool-x86_64.AppImage") {
		result = append(result, "✅ appimagetool found")
	} else {
		result = append(result, "❌ appimagetool missing")
		allPassed = false
	}

	return result, allPassed
}
