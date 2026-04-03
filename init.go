package main

import (
	"fmt"
	"os"
)

func runInit() {
	fmt.Println("📦 Creating project...")

	os.WriteFile("main.go", []byte(`package main

func main() {
	println("Hello AppImage 🚀")
}
`), 0644)

	os.WriteFile("icon.png", []byte{}, 0644)

	fmt.Println("✅ Project created")
}
