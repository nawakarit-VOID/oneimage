package main

import (
	"bytes"
	"embed"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// TEMPLATE (ฝังใน)
const buildTemplate = `#!/bin/bash
set -e
export PATH=/usr/local/go/bin:$PATH

APP={{.AppName}}
EXEC={{.ExecName}}

echo "🔍 Checking..."

#command -v go >/dev/null 2>&1 || { echo "❌ Go not found"; exit 1; }

echo "ตรวจเช็คไฟล์"
sleep 1
[ -f "icon.png" ] || { echo "❌ icon.png missing"; exit 1; }
[ -f "main.go" ] || { echo "❌ main.go missing"; exit 1; }
[ -f "go.mod" ] || { echo "❌ go.mod missing"; exit 1; }
[ -f "go.sum" ] || { echo "❌ go.sum, missing"; exit 1; }

echo "🔨 ใจเย็นๆ..."
sleep 1
go mod tidy
go build -ldflags="-s -w" -o $EXEC

echo "📦 รวมเอกสาร...AppDir..."
sleep 1
rm -rf $APP.AppDir
mkdir -p $APP.AppDir
cp $EXEC $APP.AppDir/

cat > $APP.AppDir/AppRun << 'EOF'
#!/bin/bash
HERE="$(dirname "$(readlink -f "$0")")"
exec "$HERE/{{.ExecName}}"
EOF

chmod +x $APP.AppDir/AppRun

cat > $APP.AppDir/$APP.desktop << EOF
[Desktop Entry]
Name={{.DisplayName}}
Exec={{.ExecName}}
Icon={{.ExecName}}
Type=Application
Categories={{.Categories}}
Terminal=false
EOF

cp icon.png $APP.AppDir/$EXEC.png
cp icon.png $APP.AppDir/.DirIcon

echo "🚀 pack..."
./appimagetool-x86_64.AppImage $APP.AppDir
sleep 2
cp $APP-x86_64.AppImage $APP.AppDir/$APP-x86_64.AppImage 

echo "📦 บีบอัด..tar..."
tar -czf $APP.tar.gz $APP.AppDir
sleep 2

echo "🧹 ลบ .AppDir..."
rm -rf $APP.AppDir

echo "✅ เสร็จแล้ว"
`

const iconsTemplate = `#!/bin/bash
set -e
export PATH=/usr/local/go/bin:$PATH

INPUT="icon.png"
OUTDIR="icons"

mkdir -p $OUTDIR

SIZES=(512 256 128 64)

for SIZE in "${SIZES[@]}"; do
  convert "$INPUT" \
    -resize ${SIZE}x${SIZE} \
    "$OUTDIR/icon-${SIZE}.png"
done

echo "✅ เสร็จแล้ว!"
`

type BuildConfig struct {
	AppName     string
	ExecName    string
	DisplayName string
	Categories  string
}

func copyAppImageTool(projectPath string) error {
	src := "./appimagetool-x86_64.AppImage"
	dst := filepath.Join(projectPath, "appimagetool-x86_64.AppImage")

	// ถ้ามีอยู่แล้ว → ไม่ต้อง copy
	if _, err := os.Stat(dst); err == nil {
		return nil
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0755)
}

// ============================================================================
// Icons gen+run
// ============================================================================
// generate icons.sh
func generateScripticons(projectPath string, cfg BuildConfig) error {
	tmpl, err := template.New("icons").Parse(iconsTemplate)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, cfg); err != nil {
		return err
	}

	scriptPath := filepath.Join(projectPath, "icons.sh")
	return os.WriteFile(scriptPath, out.Bytes(), 0755)
}

// run icons.sh
func runScripticons(projectPath string, output *widget.Entry) {

	commands := [][]string{
		{"gnome-terminal", "--", "bash", "-c", "cd '" + projectPath + "' && ./icons.sh; exec bash"},
		{"x-terminal-emulator", "-e", "bash", "-c", "cd '" + projectPath + "' && ./icons.sh; exec bash"},
		{"konsole", "-e", "bash", "-c", "cd '" + projectPath + "' && ./icons.sh; exec bash"},
		{"xfce4-terminal", "-e", "bash", "-c", "cd '" + projectPath + "' && ./icons.sh; exec bash"},
	}

	for _, c := range commands {
		cmd := exec.Command(c[0], c[1:]...)
		err := cmd.Start()
		if err == nil {
			output.SetText("🚀 opened terminal: " + c[0])
			return
		}
	}

	output.SetText("❌ no terminal found")
}

// ============================================================================
// build gen+run
// ============================================================================
// generate build.sh
func generateScriptbuild(projectPath string, cfg BuildConfig) error {
	tmpl, err := template.New("build").Parse(buildTemplate)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, cfg); err != nil {
		return err
	}

	scriptPath := filepath.Join(projectPath, "build.sh")
	return os.WriteFile(scriptPath, out.Bytes(), 0755)
}

// 🔥 run build.sh
func runScriptbuild(projectPath string, output *widget.Entry) {

	commands := [][]string{
		{"gnome-terminal", "--", "bash", "-c", "cd '" + projectPath + "' && ./build.sh; exec bash"},
		{"x-terminal-emulator", "-e", "bash", "-c", "cd '" + projectPath + "' && ./build.sh; exec bash"},
		{"konsole", "-e", "bash", "-c", "cd '" + projectPath + "' && ./build.sh; exec bash"},
		{"xfce4-terminal", "-e", "bash", "-c", "cd '" + projectPath + "' && ./build.sh; exec bash"},
	}

	for _, c := range commands {
		cmd := exec.Command(c[0], c[1:]...)
		err := cmd.Start()
		if err == nil {
			output.SetText("🚀 opened terminal: " + c[0])
			return
		}
	}

	output.SetText("❌ no terminal found")
}

// โหลด icon
func loadIcon(size int) fyne.Resource {
	var file string

	switch {
	case size >= 512:
		file = "icons/icon-512.png" ///ที่อยู่
	case size >= 256:
		file = "icons/icon-256.png"
	case size >= 128:
		file = "icons/icon-128.png"
	default:
		file = "icons/icon-64.png"
	}

	data, _ := iconFS.ReadFile(file)
	return fyne.NewStaticResource(file, data)
}

//go:embed icons/*
var iconFS embed.FS

// ─── Main ─────────────────────────────────────────────────────────────────────
func main() {

	a := app.NewWithID("com.nawakarit.oneimage")
	icon := loadIcon(64) //เอา data มาใช้
	a.SetIcon(icon)
	w := a.NewWindow("oneimage")
	w.SetIcon(icon)

	// 🔹 input
	appName := widget.NewEntry()
	appName.SetPlaceHolder("App Name (myapp)")

	execName := widget.NewEntry()
	execName.SetPlaceHolder("Executable Name (myapp)")

	displayName := widget.NewEntry()
	displayName.SetPlaceHolder("Display Name (My App)")

	categories := widget.NewEntry()
	categories.SetText("Utility;")

	catmenu := widget.NewMultiLineEntry()
	catmenu.SetText(`ประเภทโปรแกรม
	Utility; = ยูทิลิตี้ (ทั่วไป)
	Development; = การพัฒนา
	Game; = เกม
	Graphics; = กราฟิก
	Network; = เครือข่าย
	Office; = สำนักงาน
	Audio; = เสียง
	Video; = วิดีโอ
	System; = ระบบ`)

	catmenu.SetMinRowsVisible(10)
	//catmenu.Disable()

	projectPath := ""

	// 🔹 output
	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetPlaceHolder("Output...")

	scroll := container.NewScroll(output)
	scroll.SetMinSize(fyne.NewSize(600, 100))

	// 🔹 เลือก folder
	selectBtn := widget.NewButton("Select Project Folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}
			projectPath = uri.Path()
			err = copyAppImageTool(projectPath)
			if err != nil {
				output.SetText("❌ copy failed: " + err.Error())
				return
			}

			output.SetText("📁 Selected: " + projectPath + "\n✅ appimagetool ready")
		}, w)
	})

	// ============================================================================
	// ปุ่ม genIcon+run
	// ============================================================================
	// 🔹 generate script
	generateBtnIcon := widget.NewButton("Generate Script icons", func() {

		if projectPath == "" {
			output.SetText("❌ Please select project folder")
			return
		}

		cfg := BuildConfig{
			AppName:     appName.Text,
			ExecName:    execName.Text,
			DisplayName: displayName.Text,
			Categories:  categories.Text,
		}

		err := generateScripticons(projectPath, cfg)
		if err != nil {
			output.SetText(err.Error())
			return
		}

		output.SetText("✅ icons.sh created at:\n" + projectPath)

	})

	// 🔹 run icon scrip
	runBtnIcons := widget.NewButton("Build icons", func() {

		if projectPath == "" {
			output.SetText("❌ select folder first")
			return
		}

		//  run script
		go runScripticons(projectPath, output)

		output.SetText("🚀 Build started in terminal...")
	})

	// ============================================================================
	// ปุ่ม build+run
	// ============================================================================
	// 🔹 generate script
	generateBtnBuild := widget.NewButton("Generate Script Build", func() {

		if projectPath == "" {
			output.SetText("❌ Please select project folder")
			return
		}

		cfg := BuildConfig{
			AppName:     appName.Text,
			ExecName:    execName.Text,
			DisplayName: displayName.Text,
			Categories:  categories.Text,
		}

		err := generateScriptbuild(projectPath, cfg)
		if err != nil {
			output.SetText(err.Error())
			return
		}

		output.SetText("✅ build.sh created at:\n" + projectPath)

	})

	// 🔹 run build
	buildBtn := widget.NewButton("Run Build", func() {

		if projectPath == "" {
			output.SetText("❌ select folder first")
			return
		}

		//  copy appimagetool (ถ้ามีฟังก์ชันนี้)
		err := copyAppImageTool(projectPath)
		if err != nil {
			output.SetText("❌ copy failed: " + err.Error())
			return
		}

		//  run script
		go runScriptbuild(projectPath, output)

		output.SetText("🚀 Build started in terminal...")
	})
	// ============================================================================
	// ปุ่มรายละเอียด ฟังชั้น icons
	// ============================================================================
	btnuseicons := widget.NewButton("!", func() {
		useicons := widget.NewMultiLineEntry()
		useicons.SetText(`
// โหลด icon
func loadIcon(size int) fyne.Resource {
	var file string

	switch {
	case size >= 512:
		file = "icons/icon-512.png" ///ที่อยู่
	case size >= 256:
		file = "icons/icon-256.png"
	case size >= 128:
		file = "icons/icon-128.png"
	default:
		file = "icons/icon-64.png"
	}

	data, _ := iconFS.ReadFile(file)
	return fyne.NewStaticResource(file, data)
}

//go:embed icons/*
var iconFS embed.FS


---icon := loadIcon(64) //เอา data มาใช้

`)
		//useicons.Disable() // ทำให้แก้ไม่ได้ แต่ยัง select/copy ได้
		d := dialog.NewCustom(
			"ฟังชั้นเรียก icon มาใช้",
			"ปิด",
			useicons,
			w,
		)
		d.Resize(fyne.NewSize(500, 500)) // กว้าง 500 สูง 300
		d.Show()
	})

	// ============================================================================
	// ปุ่ม howto
	// ============================================================================
	btnhowto := widget.NewButton("#", func() {
		howto := widget.NewMultiLineEntry()
		howto.SetText(`oneimage
-ใช้สำหรับทำ .image (ไฟล์เดียว) 
-oneimage จำเป็นต้องมี ไฟล์ appimagetool-x86_64.AppImage อยู่ข้างๆเสมอ
-สามารถเอาตัว appimage version ที่ใหม่กว่ามาแทนได้เลย *แต่ชื่อต้องอ้างอิงด้านบน

**Go
**ใช้ได้กับภาษา Go
**Golang 
**fyne (gui)
**

เครื่องที่ใช้จะต้องมี 
-ภาษา go (golang) (ในเครื่อง)
-ImageMagick
-และอื่นๆที่ไฟล์ go.mod ต้องใช้งาน

แฟ้ม oneimage/
  ├── appimagetool-x86_64.AppImage (*ใช้รุ่นใหม่กว่าได้)
  └── oneimage_v1_0_0_0-x86_64.AppImage (ใช้งาน gui)

แฟ้ม**โปรเจคเป้าหมาย/
  ├── icon.png (ตั้งชื่อว่า icon.png) (Master*)
  ├── main.go
  ├── go.mod
  └── go.sum

การใช้งาน
**icon เป็นฟังชั้นเสริม**
1. ใ้ส่ชื่อ app ,exec ,Display (โดยส่วนมาก ใช้ชื่อเดีนวกันหมด)
2. ช่อง categories ก็มีให้เลือกตามข้อความ ด้านล่างช่องกรอก ปิดท้ายด้วย ; เสมอ เช่น x; , x;x; , x;x;x;
3. เลือกแฟ้มโปรเจค (ระบบจะก็อปไฟล์ appimagetool-x86_64.AppImage ไปวางในแฟ้มโปรเจค)
4. Generate script (จะมีไฟล์ Build.sh ขึ้นที่แฟ้มโปรเจค)
5. Run build`)
		//howto.Disable() // ทำให้แก้ไม่ได้ แต่ยัง select/copy ได้
		d := dialog.NewCustom(
			"howto",
			"ปิด",
			howto,
			w,
		)
		d.Resize(fyne.NewSize(600, 600)) // กว้าง 500 สูง 300
		d.Show()
	})

	// ============================================================================
	// ปุ่ม howto
	// ============================================================================
	btnlic := widget.NewButton("@", func() {
		LICENSE := widget.NewMultiLineEntry()
		LICENSE.SetText(`
Copyright (c) [2026] [nawakarit] [เจช์ (วัดดงหมี)] (https://github.com/nawakarit-VOID)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files... 

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND...

***************

ลิขสิทธิ์ (c) [2026] [nawakarit] [เจช์ (วัดดงหมี)] (https://github.com/nawakarit-VOID)

ขออนุญาตให้บุคคลใดก็ตามที่ได้รับสำเนาซอฟต์แวร์นี้และเอกสารประกอบที่เกี่ยวข้อง สามารถดาวน์โหลดไปใช้ได้โดยไม่มีค่าใช้จ่าย...

ข้อความแจ้งลิขสิทธิ์และข้อความแจ้งอนุญาตข้างต้นจะต้องรวมอยู่ในสำเนาทั้งหมดหรือส่วนสำคัญของซอฟต์แวร์

ซอฟต์แวร์นี้จัดให้ "ตามสภาพที่เป็นอยู่" โดยไม่มีการรับประกันใดๆ ทั้งสิ้น...
`)
		//howto.Disable() // ทำให้แก้ไม่ได้ แต่ยัง select/copy ได้
		d := dialog.NewCustom(
			"LICENSE",
			"ปิด",
			LICENSE,
			w,
		)
		d.Resize(fyne.NewSize(600, 600)) // กว้าง 500 สูง 300
		d.Show()
	})

	// ============================================================================
	// layout
	// ============================================================================
	// 🔹 layout

	w.SetContent(container.NewVBox(
		widget.NewLabel("⚙️ Config"),
		appName,
		execName,
		displayName,
		categories,
		catmenu,

		selectBtn,

		container.NewHBox(generateBtnIcon, runBtnIcons, widget.NewLabel("!! เอาไอคอน master มาวางที่แฟ้มงานก่อน"), btnuseicons, btnhowto, btnlic),

		container.NewHBox(generateBtnBuild, buildBtn),

		widget.NewLabel("📄 Output"),
		scroll,
	))

	w.Resize(fyne.NewSize(600, 500))
	w.ShowAndRun()
}
