package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

func runGUI() {
	// ============================================================================
	// App
	// ============================================================================
	a := app.NewWithID("com.nawakarit.oneimage")
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("oneimage")
	w.SetIcon(resourceIconPng)
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
	// ปุ่ม BUILD
	// ============================================================================
	buildBtn := widget.NewButton("Build", func() {

		cfg := BuildConfig{
			AppName:     appName.Text,
			ExecName:    execName.Text,
			DisplayName: displayName.Text,
			Type:        "Application",
			Categories:  "Utility;",
		}

		go func() {
			err := buildApp(cfg)

			if err != nil {
				output.SetText("❌ " + err.Error())
				return
			}

			output.SetText("✅ Build Complete!")
		}()
	})

	// ============================================================================
	// layout
	// ============================================================================

	top := container.NewVBox(
		appName,
		execName,
		displayName,
	)

	left := container.NewVBox()
	righ := container.NewVBox()
	bot := container.NewVBox(

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
