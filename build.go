package main

import (
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

var projectPath string

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "ARCH=x86_64")
	cmd.Dir = projectPath //

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildApp(cfg BuildConfig) error {

	logStep("preparing modules...")
	if cfg.Offline {
		runCmd("go", "mod", "tidy", "-e")
	} else {
		runCmd("go", "mod", "tidy")
	}

	logStep("building binary...")
	if err := runCmd("go", "build", "-ldflags=-s -w", "-o", cfg.ExecName); err != nil {
		return err
	}

	logStep("preparing AppDir...")
	appDir := projectPath + "/" + cfg.AppName + ".AppDir" //เก็บที่อยู่ ไฟล์ปลายทาง

	os.RemoveAll(appDir)      //ลบแฟ้มที่เหมือนกันออก
	os.MkdirAll(appDir, 0755) //สร้างแฟ้มใหม่

	copyFile(cfg.ExecName, filepath.Join(appDir, cfg.ExecName)) //

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

	logStep("packing AppImage...")
	file2 := "appimagetool-x86_64.AppImage"
	dst2 := filepath.Join(projectPath, filepath.Base(file2))
	copyFile(file2, dst2)

	if err := copyFile(file2, dst2); err != nil {
		return err
	}
	if err := os.Chmod(dst2, 0755); err != nil {
		return err
	}

	copyFile(projectPath+"/"+cfg.AppName, filepath.Join(appDir, cfg.ExecName)) //ก็อป lib มาที่ .appDir

	fmt.Print(3)
	if err := runCmd(dst2, appDir+"/"); err != nil {
		fmt.Print(7)
		return err

	}

	copyFile(projectPath+"/"+cfg.AppName+"-x86_64.AppImage", filepath.Join(appDir, cfg.AppName+"-x86_64.AppImage")) //ก็อป .appimage มาที่ .appDir

	tarName := cfg.AppName + ".tar.gz"
	appImage := cfg.AppName + ".AppDir"

	logStep("building binary...")
	if err := runCmd("tar", "-czf", tarName, appImage); err != nil {
		return err
	}

	//os.Remove(projectPath + "/" + cfg.AppName)
	//os.RemoveAll(projectPath + "/" + cfg.AppName + ".AppDir")
	os.RemoveAll(projectPath + "/" + "appimagetool-x86_64.AppImage")

	logStep("DONE ✅")
	return nil
}
