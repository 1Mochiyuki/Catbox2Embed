package main

import (
	"fmt"
	"os"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/1Mochiyuki/Catbox2Embed/fileupload"
)

func newUploadFileSection(container *fyne.Container, label *widget.Label) fyne.Widget {

	if label == nil {
		fmt.Println("label nil, using default label")
		label = widget.NewLabel("No file selected")
	}
	uploadBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), nil)
	cancelBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), nil)
	openFileBtn := widget.NewButtonWithIcon("", theme.FolderNewIcon(), nil)

	fileUploadWidget := fileupload.NewFileUploadWidget(container, uploadBtn, cancelBtn, openFileBtn, label)
	return fileUploadWidget
}

type Instructions struct {
	widget.BaseWidget
	InstructionsText string
	OnTapped         func()
}

func NewInstructions(instructionsText string, onTapped func()) *Instructions {
	fmt.Println("here")

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
	fmt.Println("tapped")
	i.OnTapped()
}

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
		for i, v := range newSlice {
			fmt.Printf("obj of %v @ : %v\n", reflect.TypeOf(v), i+1)
			fmt.Printf("len of main container obj before deletion: %v\n", len(mainContainer.Objects))
			mainContainer.Remove(mainContainer.Objects[len(mainContainer.Objects)-1])
			fmt.Printf("len of main container obj after deletion: %v\n", len(mainContainer.Objects))
			fmt.Println("----------------------")
		}
		mainContainer.Add(newUploadFileSection(mainContainer, nil))
	}

	hbox := container.NewAdaptiveGrid(2, uploadAllBtn, clearAllBtn)

	mainContainer.Add(hbox)
	fmt.Printf("len of mainContainer obj: %v\n", len(mainContainer.Objects))

	window.SetOnDropped(func(p fyne.Position, u []fyne.URI) {
		
		for _, v := range u {
			fileInfo, err := os.Stat(v.Path())
			if err != nil {
				panic(err)
			}
			sizeInMib := (fileInfo.Size() / 1024) / 1024
			if sizeInMib >= 200 {
				dialog.ShowError(fmt.Errorf("Size must be under 200 MiB. file was: %v MiB", sizeInMib), window)
				continue
			}
			fmt.Printf("file: %s\nsize of file: %v MiB\n", v.Path(), sizeInMib)
			

			uploadWidget := newUploadFileSection(mainContainer, widget.NewLabel(v.Path()))
			mainContainer.Add(uploadWidget)
		}
		if len(mainContainer.Objects) > 1 {
			fmt.Printf("len of main Container obj: %v\n", len(mainContainer.Objects))
			window.SetContent(mainContainer)
		}

	})

	if len(mainContainer.Objects) >= 0 {
		fmt.Println("nothing in ui, adding instructions")

		instructions := NewInstructions("Click or drag files to begin", func() {
			fileUploadSection := newUploadFileSection(mainContainer, nil)
			mainContainer.Add(fileUploadSection)
			window.SetContent(mainContainer)
		})

		window.SetContent(instructions)
	}

	window.ShowAndRun()
}
