package main

import (
	"fmt"
	"os/exec"
)

func checkCmd(name string) {
	_, err := exec.LookPath(name)
	if err != nil {
		fmt.Println("❌", name, "not found")
	} else {
		fmt.Println("✅", name)
	}
}

func runDoctor() {
	fmt.Println("🔍 Checking system...")

	checkCmd("go")
	checkCmd("fyne")
	checkCmd("wget")

	fmt.Println("Done")
}
