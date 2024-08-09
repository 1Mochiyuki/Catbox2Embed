package main

import (
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
	"github.com/1Mochiyuki/Catbox2Embed/fileupload"
)

/*
	TODO:
		add copy text button

*/

func addCopyAllLinksShortcut(window fyne.Window, con *fyne.Container) {
	ctrlE := &desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierControl,
	}
	window.Canvas().AddShortcut(ctrlE, func(shortcut fyne.Shortcut) {
		fmt.Println("copying links")
		var links string
		for _, v := range con.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&fileupload.FileUploadWidget{}) {
				uploadWidget := v.(*fileupload.FileUploadWidget)
				if uploadWidget.FileName.Label.Text == fileupload.DEFAULT_LABEL_TEXT {
					continue
				}

				links += uploadWidget.FileName.Label.Text + "\n"
				fmt.Printf("current links: %v\n", links)
			}
		}
		window.Clipboard().SetContent(links)
	})
}

func newUploadFileSection(app fyne.App, window fyne.Window, con *fyne.Container, absPath string) fyne.Widget {

	fileNameLabel := fileupload.NewFileNameLabel(absPath)

	uploadBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), nil)
	cancelBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), nil)
	openFileBtn := widget.NewButtonWithIcon("", theme.FolderNewIcon(), nil)
	copyTextBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		fmt.Println("bwomp")
		window.Clipboard().SetContent(fileNameLabel.Label.Text)
		app.SendNotification(fyne.NewNotification("Copy", fmt.Sprintf("Copied: %s successfully", fileNameLabel.Label.Text)))
		fmt.Println("should've sent notification")
	})

	fileUploadWidget := fileupload.NewFileUploadWidget(con, uploadBtn, cancelBtn, openFileBtn, copyTextBtn, fileNameLabel)
	return fileUploadWidget
}

type Instructions struct {
	widget.BaseWidget
	InstructionsText string
	OnTapped         func()
}

func NewInstructions(instructionsText string, onTapped func()) *Instructions {

	instructions := &Instructions{
		InstructionsText: instructionsText,
		OnTapped:         onTapped,
	}
	instructions.ExtendBaseWidget(instructions)
	return instructions
}

func (i *Instructions) CreateRenderer() fyne.WidgetRenderer {
	con := container.NewCenter(widget.NewLabelWithStyle(i.InstructionsText, fyne.TextAlignCenter, widget.RichTextStyleHeading.TextStyle))

	return widget.NewSimpleRenderer(con)
}
func (i *Instructions) Tapped(ev *fyne.PointEvent) {

	i.OnTapped()
}

var FILE_SIZE_REQUIREMENT = 200

func main() {

	a := app.NewWithID("Catbox2Embed")
	window := a.NewWindow("Catbox2Embed")
	window.Resize(fyne.NewSize(600, 500))
	mainContainer := container.NewVBox()

	uploadAllBtn := widget.NewButton("Upload All", func() {})
	uploadAllBtn.OnTapped = func() {
		for _, v := range mainContainer.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&fileupload.FileUploadWidget{}) {
				uploadWidget := v.(*fileupload.FileUploadWidget)

				uploadWidget.StartUploadButton.OnTapped()
			}
		}
	}
	clearAllBtn := widget.NewButton("Clear All", func() {})

	clearAllBtn.OnTapped = func() {
		newSlice := mainContainer.Objects[1:]
		fmt.Printf("len of mainContainer obj: %v\n", len(mainContainer.Objects))
		fmt.Printf("len of new slice: %v\n", len(newSlice))

		for i := 0; i < len(newSlice); i++ {
			mainContainer.Remove(mainContainer.Objects[len(mainContainer.Objects)-1])
		}
		mainContainer.Add(newUploadFileSection(a, window, mainContainer, fileupload.DEFAULT_LABEL_TEXT))
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
			if sizeInMib > int64(FILE_SIZE_REQUIREMENT) {
				dialog.ShowError(fmt.Errorf("size must be %v MiB or lower. file was: %v MiB", FILE_SIZE_REQUIREMENT, sizeInMib), window)
				continue
			}

			uploadWidget := newUploadFileSection(a, window, mainContainer, v.Path())
			mainContainer.Add(uploadWidget)
		}
		if len(mainContainer.Objects) > 1 {
			addCopyAllLinksShortcut(window, mainContainer)
			window.SetContent(mainContainer)
		}

	})

	if len(mainContainer.Objects) >= 0 {

		instructions := NewInstructions(fmt.Sprintf("Click or drag files to begin. Must be %v MiB or lower\nPress Ctrl + E while window is focused to copy all links", FILE_SIZE_REQUIREMENT), func() {
			fileUploadSection := newUploadFileSection(a, window, mainContainer, fileupload.DEFAULT_LABEL_TEXT)
			mainContainer.Add(fileUploadSection)
			window.SetContent(mainContainer)
		})

		window.SetContent(instructions)
	}

	window.ShowAndRun()
}
