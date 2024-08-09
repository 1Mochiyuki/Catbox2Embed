package fileupload

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"github.com/wabarc/go-catbox"
)

type FileNameLabel struct {
	AbsPath string
	Label   *widget.Label
}

var DEFAULT_LABEL_TEXT = "No file selected"

func NewFileNameLabel(absPath string) *FileNameLabel {
	label := widget.NewLabel(absPath)

	if strings.Contains(absPath, "/") {
		fmt.Println("meow")
		parts := strings.Split(absPath, "/")
		fileName := parts[len(parts)-1]
		fmt.Printf("filename: %v\n", absPath)
		label.SetText(fileName)
	}

	return &FileNameLabel{
		AbsPath: absPath,
		Label:   label,
	}
}

type FileUploadWidget struct {
	widget.BaseWidget
	OpenFileButton     *widget.Button
	StartUploadButton  *widget.Button
	CancelUploadButton *widget.Button
	FileName           *FileNameLabel
}

func uploadToCatbox(fileName *FileNameLabel) {
	if fileName.Label.Text == DEFAULT_LABEL_TEXT {
		fmt.Println("no file selected, not uploading")
		return
	}
	filePath := fileName.AbsPath
	fmt.Printf("uploading: %s\n", filePath)
	if url, err := catbox.New(nil).Upload(filePath); err != nil {
		fmt.Printf("\ncatbox: %v\n", err)
		return
	} else {
		embedLink := fmt.Sprintf("https://embeds.video/%s", url)
		fmt.Printf("\nurl: %s\npath: %s", embedLink, filePath)
		fileName.Label.SetText(embedLink)
	}
}

func NewFileUploadWidget(container *fyne.Container, startUploadButton, cancelUploadButton, openFileButton *widget.Button, fileName *FileNameLabel) *FileUploadWidget {

	item := &FileUploadWidget{
		OpenFileButton:     openFileButton,
		StartUploadButton:  startUploadButton,
		CancelUploadButton: cancelUploadButton,
		FileName:           fileName,
	}

	//item.FileName.Label.Alignment = fyne.TextAlignCenter

	cancelUploadButton.OnTapped = func() {

		fileName.Label.SetText(DEFAULT_LABEL_TEXT)
		if len(container.Objects) > 1 {
			fmt.Printf("length of container objects: %v\n", len(container.Objects))
			objPlace := len(container.Objects) - 1
			container.Remove(container.Objects[objPlace])
			return
		}
		fmt.Println("file name cleared, nothing else happened")
	}

	startUploadButton.OnTapped = func() {
		go uploadToCatbox(fileName)
	}

	openFileButton.OnTapped = func() {
		fmt.Println("set open file tapped func")
		filename, err := dialog.File().Load()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fileName.Label.SetText(filename)

		nextUploadBtn := widget.NewButtonWithIcon(startUploadButton.Text, startUploadButton.Icon, startUploadButton.OnTapped)
		nextCancelBtn := widget.NewButtonWithIcon(cancelUploadButton.Text, cancelUploadButton.Icon, cancelUploadButton.OnTapped)
		nextOpenFileBtn := widget.NewButtonWithIcon(openFileButton.Text, openFileButton.Icon, openFileButton.OnTapped)

		container.Add(NewFileUploadWidget(container, nextUploadBtn, nextCancelBtn, nextOpenFileBtn, NewFileNameLabel(DEFAULT_LABEL_TEXT)))
	}
	item.ExtendBaseWidget(item)
	return item
}

func (item *FileUploadWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.CancelUploadButton, item.StartUploadButton, item.OpenFileButton, item.FileName.Label)
	return widget.NewSimpleRenderer(c)
}
