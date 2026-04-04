package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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

//go:embed icon.png
var iconData []byte

func runGUI() {
	// ============================================================================
	// App
	// ============================================================================

	a := app.NewWithID("com.nawakarit.oneimage")
	icon := fyne.NewStaticResource("icon.png", iconData)
	a.SetIcon(icon)
	w := a.NewWindow("oneimage")
	w.SetIcon(icon)

	// ============================================================================
	// กล่องใส่ข้อความ + เลือก Type + เลือก Categories
	// ============================================================================
	appName := widget.NewEntry()
	appName.SetPlaceHolder("App Name (myapp)")
	execName := widget.NewEntry()
	execName.SetPlaceHolder("Executable Name (myapp)")
	displayName := widget.NewEntry()
	displayName.SetPlaceHolder("displayName Name (myapp)")

	// ============================================================================
	// Output
	// ============================================================================
	output := widget.NewMultiLineEntry()

	// ============================================================================
	// tyne
	// ============================================================================

	typeSelect := widget.NewSelect([]string{
		"Application",
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
	var categoryChecks []*widget.Check

	for _, c := range categoriesList {
		name := c
		check := widget.NewCheck(name, func(bool) {})
		categoryChecks = append(categoryChecks, check)
	}

	if len(categoryChecks) > 0 {
		categoryChecks[0].SetChecked(true) // Utility
	}

	var categoryObjects []fyne.CanvasObject
	for _, c := range categoryChecks {
		categoryObjects = append(categoryObjects, c)
	}

	// ============================================================================
	// ปุ่ม BUILD
	// ============================================================================
	//var buildBtn *widget.Button
	//buildBtn.Disable()
	buildBtn := widget.NewButton("Build", func() {

		cfg := BuildConfig{
			AppName:     appName.Text,
			ExecName:    execName.Text,
			DisplayName: displayName.Text,
			Type:        typeSelect.Selected,
			Categories:  getCategories(categoryChecks),
		}

		go func() {
			err := buildApp(cfg)

			fyne.Do(func() {
				if err != nil {
					output.SetText("❌ " + err.Error())
					return
				}
				output.SetText("✅ Build Complete!")
			})
		}()
	})

	// ============================================================================
	// เลือกแฟ้ม
	// ============================================================================

	selectFolderBtn := widget.NewButton("📂 Select Project", func() {

		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}

			projectPath = uri.Path()
			output.SetText("📂 Selected: " + projectPath)

		}, w).Show()

	})

	// ============================================================================
	// ปุ่ม เช็ค list
	// ============================================================================
	checkBtn := widget.NewButton("🔍 Check System", func() {

		if projectPath == "" {
			output.SetText("❌ กรุณาเลือกโฟลเดอร์ก่อน")
			return
		}

		results, allPassed := runDoctor(projectPath)

		output.SetText("")
		for _, r := range results {
			output.SetText(output.Text + r + "\n")
		}

		if allPassed {
			buildBtn.Enable()
		} else {
			buildBtn.Disable()
		}
	})

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

	// ============================================================================
	// layout
	// ============================================================================

	top := container.NewVBox(
		appName,
		execName,
		displayName,
		widget.NewLabel("Type"),
		typeSelect,
		widget.NewLabel("Categories"),
		container.NewVBox(categoryObjects...),
	)

	left := container.NewVBox()
	righ := container.NewVBox()
	bot := container.NewVBox(
		selectAllBtn,
		clearBtn,
		selectFolderBtn,
		checkBtn,
		buildBtn,
	)
	//
	w.SetContent(container.NewVBox(
		top,
		left,
		righ,
		bot,
		output,
	))

	w.Resize(fyne.NewSize(500, 500))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
