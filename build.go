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

func zipFolder(source, target string) error {

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// สร้าง path ใน zip
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// ถ้าเป็น folder → skip (แต่ต้องมี /)
		if info.IsDir() {
			if relPath == "." {
				return nil
			}
			_, err := archive.Create(relPath + "/")
			return err
		}

		// เปิดไฟล์
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// สร้าง file ใน zip
		writer, err := archive.Create(relPath)
		if err != nil {
			return err
		}

		// copy data
		_, err = io.Copy(writer, file)
		return err
	})
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

	thisfile1 := appDir
	dst1 := filepath.Join(projectPath, filepath.Base(thisfile1))
	copyFile(thisfile1, dst1)
	fmt.Print(1)

	logStep("packing AppImage...")
	thisfile2 := "appimagetool-x86_64.AppImage"
	dst2 := filepath.Join(projectPath, filepath.Base(thisfile2))
	copyFile(thisfile2, dst2)

	if err := copyFile(thisfile2, dst2); err != nil {
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

	copyFile(projectPath+"/"+cfg.AppName+"-x86_64.AppImage", filepath.Join(appDir, cfg.ExecName+"-x86_64.AppImage")) //ก็อป lib มาที่ .appDir
	//-x86_64.AppImage

	logStep("Zip")
	err := zipFolder(appDir, appDir+cfg.AppName+".zip") //"/"
	if err != nil {
		fmt.Print("❌ zip fail: " + err.Error())
		return err
	}

	os.Remove(projectPath + "/" + cfg.AppName)
	os.RemoveAll(projectPath + "/" + cfg.AppName + ".AppDir")
	os.RemoveAll(projectPath + "/" + "appimagetool-x86_64.AppImage")

	logStep("DONE ✅")
	return nil
}
