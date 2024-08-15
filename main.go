package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/1Mochiyuki/Catbox2Embed/ui"
	"github.com/1Mochiyuki/Catbox2Embed/utils"
)

func addCopyAllLinksShortcut(window fyne.Window, con *fyne.Container) {
	ctrlE := &desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierControl,
	}
	window.Canvas().AddShortcut(ctrlE, func(shortcut fyne.Shortcut) {
		fmt.Println("copying links")
		var links string
		var linksCount int
		for _, v := range con.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&ui.FileUploadWidget{}) {
				uploadWidget := v.(*ui.FileUploadWidget)
				if uploadWidget.FileName.Label.Text == ui.DEFAULT_LABEL_TEXT {
					continue
				}

				links += uploadWidget.FileName.Label.Text + "\n"
				linksCount += 1
				fmt.Printf("current links: %v\n", links)
			}
		}
		window.Clipboard().SetContent(links)
		if utils.PreferencesEnabled() {

			fyne.CurrentApp().SendNotification(fyne.NewNotification("Copy All", fmt.Sprintf("Copied %v links", linksCount)))
		}
		links = ""
		linksCount = 0
	})
}

func main() {

	a := app.NewWithID("Catbox2Embed")
	window := a.NewWindow("Catbox2Embed")
	window.Resize(fyne.NewSize(600, 500))
	mainContainer := container.NewVBox()

	copyAllToolbarAction := widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		var links string
		for _, v := range mainContainer.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&ui.FileUploadWidget{}) {
				uploadWidget := v.(*ui.FileUploadWidget)
				if uploadWidget.FileName.Label.Text == ui.DEFAULT_LABEL_TEXT {
					continue
				}

				links += uploadWidget.FileName.Label.Text + "\n"
				fmt.Printf("current links: %v\n", links)
			}
		}
		window.Clipboard().SetContent(links)
	})

	toolbar := ui.CreateToolBar(mainContainer, a, window, copyAllToolbarAction)

	content := container.NewBorder(toolbar, nil, nil, nil)
	mainContainer.Add(content)

	uploadAllBtn := widget.NewButton("Upload All", func() {})
	uploadAllBtn.OnTapped = func() {
		for _, v := range mainContainer.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&ui.FileUploadWidget{}) {
				uploadWidget := v.(*ui.FileUploadWidget)

				uploadWidget.StartUploadButton.OnTapped()
				if copyAllToolbarAction.Disabled() {
					copyAllToolbarAction.Enable()
				}
			}
		}
	}
	clearAllBtn := widget.NewButton("Clear All", func() {})

	clearAllBtn.OnTapped = func() {
		newSlice := mainContainer.Objects[3:]
		fmt.Printf("len of mainContainer obj: %v\n", len(mainContainer.Objects))
		fmt.Printf("len of new slice: %v\n", len(newSlice))

		for i := 0; i < len(newSlice); i++ {
			mainContainer.Remove(mainContainer.Objects[len(mainContainer.Objects)-1])
		}
		if !copyAllToolbarAction.Disabled() {
			copyAllToolbarAction.Disable()
		}

	}

	hbox := container.NewAdaptiveGrid(2, uploadAllBtn, clearAllBtn)

	mainContainer.Add(hbox)

	window.SetOnDropped(func(p fyne.Position, u []fyne.URI) {

		for _, v := range u {
			fileInfo, err := os.Stat(v.Path())
			if err != nil {
				panic(err)
			}
			sizeInMib := (fileInfo.Size() / 1024) / 1024
			if sizeInMib > int64(utils.FILE_SIZE_REQUIREMENT) {
				dialog.ShowError(fmt.Errorf("size must be %v MiB or lower. file was: %v MiB", utils.FILE_SIZE_REQUIREMENT, sizeInMib), window)
				continue
			}
			if !utils.IsVideoFile(fileInfo.Name()) {
				dialog.ShowError(errors.New("file is not a video.\nif this is an error, please contact the developer"), window)
				continue
			}

			uploadWidget := ui.NewUploadFileSection(a, window, mainContainer, v.Path())
			mainContainer.Add(uploadWidget)
		}
		if len(mainContainer.Objects) > 1 {
			addCopyAllLinksShortcut(window, mainContainer)
			scrollContainer := container.NewVScroll(mainContainer)
			window.SetContent(scrollContainer)
		}
	})

	if len(mainContainer.Objects) >= 0 {

		instructions := ui.NewInstructions(fmt.Sprintf("Click or drag files to begin. Must be %v MiB or lower\nPress Ctrl + E while window is focused to copy all links", utils.FILE_SIZE_REQUIREMENT), func() {
			fileUploadSection := ui.NewUploadFileSection(a, window, mainContainer, ui.DEFAULT_LABEL_TEXT)
			mainContainer.Add(fileUploadSection)
			scrollContainer := container.NewVScroll(mainContainer)
			window.SetContent(scrollContainer)
		})

		window.SetContent(instructions)
	}
	window.SetOnClosed(func() {
		a.Quit()
	})
	window.ShowAndRun()
}
