package main

import (
	"bytes"
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

// 🔥 TEMPLATE (ฝังใน Go เลย)
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
[ -f "go.mod" ] || { echo "❌ main.go missing"; exit 1; }
[ -f "go.sum" ] || { echo "❌ main.go missing"; exit 1; }

echo "🔨 build..."
sleep 1
go mod tidy
go build -ldflags="-s -w" -o $EXEC

echo "📦 prepare...AppDir..."
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

echo "📦 tar..."
tar -czf $APP.tar.gz $APP.AppDir
sleep 2

echo "🧹 cleanup..."
rm -rf $APP.AppDir

echo "✅ DONE"
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

// 🔥 generate build.sh
func generateScript(projectPath string, cfg BuildConfig) error {
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
func runScript(projectPath string, output *widget.Entry) {

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

//go:embed icons/icon-64.png
var iconData []byte

// ─── Main ─────────────────────────────────────────────────────────────────────
func main() {

	a := app.NewWithID("com.nawakarit.oneimage")
	icon := fyne.NewStaticResource("icon-64.png", iconData)
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
	catmenu.SetText("ประเภทโปรแกรม Utility;Development;Game;Graphics;Network;Office;AudioVideo;System;")
	//catmenu.Disable()

	projectPath := ""

	// 🔹 output
	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetPlaceHolder("Output...")

	scroll := container.NewScroll(output)
	scroll.SetMinSize(fyne.NewSize(500, 300))

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

	// 🔹 generate script
	generateBtn := widget.NewButton("Generate Script", func() {

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

		err := generateScript(projectPath, cfg)
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
		go runScript(projectPath, output)

		output.SetText("🚀 Build started in terminal...")
	})

	// 🔹 layout
	w.SetContent(container.NewVBox(
		widget.NewLabel("⚙️ Config"),

		appName,
		execName,
		displayName,
		categories,
		catmenu,

		selectBtn,

		container.NewHBox(generateBtn, buildBtn),

		widget.NewLabel("📄 Output"),
		scroll,
	))

	w.Resize(fyne.NewSize(600, 500))
	w.ShowAndRun()
}
