package main

import (
	"archive/zip"
	"fmt"
	"io"
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

var projectPath string

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "ARCH=x86_64")
	cmd.Dir = projectPath //

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

/*func checkFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("❌ missing file: %s", path)
	}
	return nil
}*/

/*func ensureFyne(offline bool) error {
	_, err := exec.LookPath("fyne")
	if err == nil {
		return nil
	}

	if offline {
		return errors.New("❌ fyne CLI not found (offline mode)")
	}

	logStep("installing fyne CLI...")
	return runCmd("go", "install", "fyne.io/fyne/v2/cmd/fyne@latest")
}*/

/*func ensureAppImageTool(offline bool) error {
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
}*/

func buildApp(cfg BuildConfig) error {

	/*logStep("checking required files...")

	if err := checkFile("icon.png"); err != nil {
		return err
	}
	if err := checkFile("main.go"); err != nil {
		return err
	}*/

	/*if err := ensureFyne(cfg.Offline); err != nil {
		return err
	}*/
	/*if err := ensureAppImageTool(cfg.Offline); err != nil {
		return err
	}*/

	logStep("preparing modules...")
	if cfg.Offline {
		runCmd("go", "mod", "tidy", "-e")
	} else {
		runCmd("go", "mod", "tidy")
	}

	/*logStep("bundling icon...")
	os.Remove("bundled.go")

	out, err := os.Create("bundled.go")
	if err != nil {
		return err
	}
	defer out.Close()*/

	/*cmd := exec.Command("fyne", "bundle", "icon.png")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	*/
	logStep("building binary...")
	if err := runCmd("go", "build", "-ldflags=-s -w", "-o", cfg.ExecName); err != nil {
		return err
	}

	logStep("preparing AppDir...")
	appDir := projectPath + "/" + cfg.AppName + ".AppDir"

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
	copyFile(projectPath+"/"+"icon.png", filepath.Join(appDir, cfg.ExecName+".png"))
	copyFile(projectPath+"/"+"icon.png", filepath.Join(appDir, ".DirIcon"))

	iconPath := filepath.Join(appDir, "usr/share/icons/hicolor/256x256/apps")
	os.MkdirAll(iconPath, 0755)
	copyFile(projectPath+"/"+"icon.png", filepath.Join(iconPath, cfg.ExecName+".png"))

	//copyFile(appDir, filepath.Join(projectPath+"/"))

	src2 := appDir
	dst2 := filepath.Join(projectPath, filepath.Base(src2))

	fmt.Print(1)
	copyFile(src2, dst2)

	//os.Remove(cfg.AppName + ".AppImage")

	logStep("packing AppImage...")
	src := "appimagetool-x86_64.AppImage"
	dst := filepath.Join(projectPath, filepath.Base(src))

	fmt.Print(1)
	copyFile(src, dst)
	fmt.Print(2)
	if err := copyFile(src, dst); err != nil {
		return err

	}
	fmt.Print(3)
	if err := os.Chmod(dst, 0755); err != nil {
		return err

	}

	copyFile(projectPath+"/"+cfg.AppName, filepath.Join(appDir, cfg.ExecName)) ///

	fmt.Print(3)
	if err := runCmd(dst, appDir+"/"); err != nil {
		fmt.Print(7)
		return err

	}

	/*fmt.Print(4)
	if err := os.Chmod("appimagetool-x86_64.AppImage", 0755); err != nil {
		fmt.Print(5)
		return err
	}*/
	/*
		if err := runCmd("./appimagetool-x86_64.AppImage", appDir+"/"); err != nil {
			fmt.Print(6)
			return err
		}*/

	//os.Remove(cfg.AppName + ".AppImage")
	//os.RemoveAll(cfg.AppName + ".AppDir")
	//os.Rename(cfg.AppName+".AppImage.zip", projectPath+"/myapp.zip")

	logStep("DONE ✅")
	return nil
}

func zipFile(src, dst string) error {
	zipfile, _ := os.Create(dst)
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	file, _ := os.Open(src)
	defer file.Close()

	w, _ := archive.Create(filepath.Base(src))
	_, err := io.Copy(w, file)
	return err
}
