package fileupload

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"github.com/wabarc/go-catbox"
)

type FileUploadWidget struct {
	widget.BaseWidget
	OpenFileButton     *widget.Button
	StartUploadButton  *widget.Button
	CancelUploadButton *widget.Button
	FileName           *widget.Label
}

var DEFAULT_LABEL_TEXT = "No file selected"

func NewFileUploadWidget(container *fyne.Container, startUploadButton, cancelUploadButton, openFileButton *widget.Button, fileName *widget.Label) *FileUploadWidget {

	item := &FileUploadWidget{
		OpenFileButton:     openFileButton,
		StartUploadButton:  startUploadButton,
		CancelUploadButton: cancelUploadButton,
		FileName:           fileName,
	}

	item.FileName.Alignment = fyne.TextAlignCenter

	cancelUploadButton.OnTapped = func() {
		fileName.SetText(DEFAULT_LABEL_TEXT)
		if len(container.Objects) > 1 { 
		fmt.Printf("length of container objects: %v\n", len(container.Objects))
		objPlace := len(container.Objects) - 1
			container.Remove(container.Objects[objPlace])
			return
		}
		fmt.Println("file name cleared, nothing else happened")
	}

	startUploadButton.OnTapped = func() {
		if fileName.Text == DEFAULT_LABEL_TEXT {
			fmt.Println("no file selected, not uploading")
			return
		}
		filePath := fileName.Text
		fmt.Printf("uploading: %s\n", filePath)
		if url, err := catbox.New(nil).Upload(filePath); err != nil {
			fmt.Printf("\ncatbox: %v\n", err)
			return
		} else {
			fmt.Printf("\nurl: https://embeds.video/%s\npath: %s", url, filePath)
		}
	}

	openFileButton.OnTapped = func() {
		fmt.Println("set open file tapped func")
		filename, err := dialog.File().Load()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fileName.SetText(filename)

		nextUploadBtn := widget.NewButtonWithIcon(startUploadButton.Text, startUploadButton.Icon, startUploadButton.OnTapped)
		nextCancelBtn := widget.NewButtonWithIcon(cancelUploadButton.Text, cancelUploadButton.Icon, cancelUploadButton.OnTapped)
		nextOpenFileBtn := widget.NewButtonWithIcon(openFileButton.Text, openFileButton.Icon, openFileButton.OnTapped)

		container.Add(NewFileUploadWidget(container, nextUploadBtn, nextCancelBtn, nextOpenFileBtn, widget.NewLabel(DEFAULT_LABEL_TEXT)))
	}
	item.ExtendBaseWidget(item)
	return item
}

func (item *FileUploadWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.CancelUploadButton, item.StartUploadButton, item.OpenFileButton, item.FileName)
	return widget.NewSimpleRenderer(c)
}
