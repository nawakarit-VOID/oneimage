package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type BuildConfig struct {
	AppName     string
	ExecName    string
	DisplayName string
	Type        string
	Categories  string
	Offline     bool
}

func logStep(msg string) {
	fmt.Println("🔹", msg)
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("❌ missing file: %s", path)
	}
	return nil
}

func ensureFyne(offline bool) error {
	_, err := exec.LookPath("fyne")
	if err == nil {
		return nil
	}

	if offline {
		return errors.New("❌ fyne CLI not found (offline mode)")
	}

	logStep("installing fyne CLI...")
	return runCmd("go", "install", "fyne.io/fyne/v2/cmd/fyne@latest")
}

func ensureAppImageTool(offline bool) error {
	if _, err := os.Stat("appimagetool-x86_64.AppImage"); err == nil {
		return nil
	}

	if offline {
		return errors.New("❌ appimagetool missing (offline mode)")
	}

	logStep("downloading appimagetool...")
	return runCmd("wget",
		"https://github.com/AppImage/AppImageKit/releases/latest/download/appimagetool-x86_64.AppImage",
	)
}

func buildApp(cfg BuildConfig) error {

	logStep("checking required files...")

	if err := checkFile("icon.png"); err != nil {
		return err
	}
	if err := checkFile("main.go"); err != nil {
		return err
	}

	if err := ensureFyne(cfg.Offline); err != nil {
		return err
	}
	if err := ensureAppImageTool(cfg.Offline); err != nil {
		return err
	}

	logStep("preparing modules...")
	if cfg.Offline {
		runCmd("go", "mod", "tidy", "-e")
	} else {
		runCmd("go", "mod", "tidy")
	}

	logStep("bundling icon...")
	os.Remove("bundled.go")

	out, err := os.Create("bundled.go")
	if err != nil {
		return err
	}
	defer out.Close()

	cmd := exec.Command("fyne", "bundle", "icon.png")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	logStep("building binary...")
	if err := runCmd("go", "build", "-ldflags=-s -w", "-o", cfg.ExecName); err != nil {
		return err
	}

	logStep("preparing AppDir...")
	appDir := cfg.AppName + ".AppDir"

	os.RemoveAll(appDir)
	os.MkdirAll(appDir, 0755)

	copyFile(cfg.ExecName, filepath.Join(appDir, cfg.ExecName))

	// AppRun
	appRun := fmt.Sprintf(`#!/bin/bash
HERE="$(dirname "$(readlink -f "$0")")"
exec "$HERE/%s"
`, cfg.ExecName)

	os.WriteFile(filepath.Join(appDir, "AppRun"), []byte(appRun), 0755)

	// Desktop
	desktop := fmt.Sprintf(`[Desktop Entry]
Name=%s
Exec=%s
Icon=%s
Type=%s
Categories=%s
Terminal=false
`, cfg.DisplayName, cfg.ExecName, cfg.ExecName, cfg.Type, cfg.Categories)

	os.WriteFile(filepath.Join(appDir, cfg.AppName+".desktop"), []byte(desktop), 0644)

	// Icons
	copyFile("icon.png", filepath.Join(appDir, cfg.ExecName+".png"))
	copyFile("icon.png", filepath.Join(appDir, ".DirIcon"))

	iconPath := filepath.Join(appDir, "usr/share/icons/hicolor/256x256/apps")
	os.MkdirAll(iconPath, 0755)
	copyFile("icon.png", filepath.Join(iconPath, cfg.ExecName+".png"))

	logStep("packing AppImage...")

	if err := os.Chmod("appimagetool-x86_64.AppImage", 0755); err != nil {
		return err
	}

	if err := runCmd("./appimagetool-x86_64.AppImage", appDir); err != nil {
		return err
	}

	logStep("DONE ✅")
	return nil
}
