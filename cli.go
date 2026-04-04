package main

import (
	"fmt"
	"os"
)

func runCLI() {
	cmd := os.Args[1]

	switch cmd {

	case "build":
		cfg := BuildConfig{
			AppName:     "myapp",
			ExecName:    "myapp",
			DisplayName: "My App",
			Type:        "Application",
			Categories:  "Utility;",
		}

		err := buildApp(cfg)
		if err != nil {
			fmt.Println("❌", err)
			os.Exit(1)
		}

	case "doctor":
		//		runDoctor()

	case "init":
		runInit()

	default:
		fmt.Println("❌ unknown command")
	}
}
