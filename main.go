package main

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type BuildConfig struct {
	AppName     string
	ExecName    string
	DisplayName string
}

func generateScript(cfg BuildConfig) (string, error) {
	data, err := os.ReadFile("templates/build.sh.tpl")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("build").Parse(string(data))
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, cfg)
	if err != nil {
		return "", err
	}

	err = os.WriteFile("build.sh", out.Bytes(), 0755)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func runBuild(output *widget.Entry) {
	cmd := exec.Command("bash", "build.sh")

	out, err := cmd.CombinedOutput()
	if err != nil {
		output.SetText(string(out) + "\nERROR: " + err.Error())
		return
	}

	output.SetText(string(out))
}

func main() {

	a := app.NewWithID("com.nawakarit.oneimage")
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("oneimage")
	w.SetIcon(resourceIconPng)

	// inputs
	appName := widget.NewEntry()
	appName.SetPlaceHolder("App Name (myapp)")

	execName := widget.NewEntry()
	execName.SetPlaceHolder("Executable Name (myapp)")

	displayName := widget.NewEntry()
	displayName.SetPlaceHolder("Display Name (My App)")

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Output...")
	output.Wrapping = fyne.TextWrapWord

	// buttons
	generateBtn := widget.NewButton("Generate Script", func() {
		cfg := BuildConfig{
			AppName:     appName.Text,
			ExecName:    execName.Text,
			DisplayName: displayName.Text,
		}

		script, err := generateScript(cfg)
		if err != nil {
			output.SetText(err.Error())
			return
		}

		output.SetText("✅ Script Generated:\n\n" + script)
	})

	buildBtn := widget.NewButton("Build AppImage", func() {
		runBuild(output)
	})

	// layout
	content := container.NewVBox(
		widget.NewLabel("⚙️ Config"),
		appName,
		execName,
		displayName,
		container.NewHBox(generateBtn, buildBtn),
		widget.NewLabel("📄 Output"),
		output,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 500))
	w.ShowAndRun()
}
