package main

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type BuildConfig struct {
	AppName     string
	ExecName    string
	DisplayName string
	Type        string
	Categories  string
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

var isGenerated bool //

// Categories เป็น string
func getCategories(checks []*widget.Check) string {
	var result string
	for _, c := range checks {
		if c.Checked {
			result += c.Text + ";"
		}
	}
	return result
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

	status := widget.NewLabel("Ready")

	typeSelect := widget.NewSelect([]string{
		"Application",
		"Link",
		"Directory",
	}, func(value string) {})

	typeSelect.SetSelected("Application")

	categoriesList := []string{
		"Utility",
		"Development",
		"Game",
		"Graphics",
		"Network",
		"Office",
		"AudioVideo",
		"System",
	}
	//categoryChecks
	var categoryChecks []*widget.Check

	for _, c := range categoriesList {

		check := widget.NewCheck(c, nil)
		categoryChecks = append(categoryChecks, check)
	}

	//categoryChecks set default
	if len(categoryChecks) > 0 {
		categoryChecks[0].SetChecked(true) // Utility
	}

	var categoryObjects []fyne.CanvasObject

	for _, c := range categoryChecks {
		categoryObjects = append(categoryObjects, c)
	}

	categoryBox := container.NewVBox(categoryObjects...)
	categoryScroll := container.NewVScroll(categoryBox)
	categoryScroll.SetMinSize(fyne.NewSize(200, 350)) //scoll

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Output...")
	output.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(output)
	scroll.SetMinSize(fyne.NewSize(350, 600))

	// buttons set
	generateBtn := widget.NewButton("Generate Script", nil)
	buildBtn := widget.NewButton("Build AppImage", nil)

	generateBtn.Disable()
	buildBtn.Disable()

	updateState := func() {
		// 🔹 validate input
		if appName.Text != "" && execName.Text != "" && displayName.Text != "" {
			generateBtn.Enable()
		} else {
			generateBtn.Disable()
		}

		// 🔥 reset build state ทุกครั้งที่มีการแก้
		isGenerated = false
		buildBtn.Disable()
	}

	appName.OnChanged = func(string) {
		updateState()
	}
	execName.OnChanged = func(string) {
		updateState()
	}
	displayName.OnChanged = func(string) {
		updateState()
	}

	generateBtn.OnTapped = func() {
		status.SetText("⚙️ Generating script...")

		cfg := BuildConfig{
			AppName:     appName.Text,
			ExecName:    execName.Text,
			DisplayName: displayName.Text,
			Type:        typeSelect.Selected,
			Categories:  getCategories(categoryChecks),
		}

		script, err := generateScript(cfg)
		if err != nil {
			output.SetText(err.Error())
			status.SetText("🔴 Generate failed")
			return
		}

		output.SetText("✅ Script Generated:\n\n" + script)
		status.SetText("✅ Script generated")

		// เปิดปุ่ม build
		isGenerated = true
		buildBtn.Enable()
		////////text setup/////////
		popup := a.NewWindow("📌 Setup Guide")

		text := widget.NewMultiLineEntry()
		text.SetText(`
		//เอาไปวางในตำแหน่ง func main()
		
		//func main() {
				a := app.NewWithID("com.yourname.yourapp")
				a.SetIcon(resourceIconPng)

				w := a.NewWindow("yourApp")
				w.SetIcon(resourceIconPng)
				
				//..............................your code................................}
				
				`)

		text.Wrapping = fyne.TextWrapWord
		text.Disable()

		scroll := container.NewScroll(text)
		scroll.SetMinSize(fyne.NewSize(600, 300)) //ยืดได้

		popup.SetContent(container.NewVBox(
			widget.NewLabel("ก่อนกด Build AppImage ให้ Copy โค้ดด้านล่าง 👇"),
			scroll,
			widget.NewButton("Copy", func() {
				a.Clipboard().SetContent(text.Text)

			}),
		))

		popup.Resize(fyne.NewSize(600, 400))
		popup.SetFixedSize(true)
		popup.Show()

	}

	buildBtn.OnTapped = func() {
		if !isGenerated {
			status.SetText("❌ Please generate script first")
			return
		}

		status.SetText("🚀 Building AppImage...")

		runBuild(output)

		status.SetText("✅ Build finished")
	}

	selectAllBtn := widget.NewButton("Select All", func() {
		for _, c := range categoryChecks {
			c.SetChecked(true)
		}
	})

	clearBtn := widget.NewButton("Clear", func() {
		for _, c := range categoryChecks {
			c.SetChecked(false)
		}
	})

	// status
	statusBar := container.NewBorder(
		nil, nil, nil, nil,
		container.NewPadded(status),
	)
	// layout/////

	card := container.NewGridWithColumns(2,

		// 🔹 LEFT PANEL
		container.NewVBox(
			widget.NewLabel("⚙️ Config"),
			appName,
			execName,
			displayName,

			layout.NewSpacer(), // ดันปุ่มไปล่าง

			widget.NewLabel("Type"),
			typeSelect,

			widget.NewLabel("Categories"),
			//container.NewVBox(categoryObjects...),
			categoryScroll,

			container.NewHBox(selectAllBtn, clearBtn),

			layout.NewSpacer(), // ดันปุ่มไปล่าง

			container.NewHBox(generateBtn, buildBtn),
		),

		// 🔹 RIGHT PANEL
		container.NewBorder(
			widget.NewLabel("📄 Output"),
			nil,
			nil,
			nil,
			scroll,
		),
	)

	w.SetContent(container.NewBorder(
		nil,       // top
		statusBar, // bottom
		nil,
		nil,
		card, // center
	))

	w.Resize(fyne.NewSize(900, 600))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
