package ui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/1Mochiyuki/Catbox2Embed/utils"
	"github.com/sqweek/dialog"
	"github.com/wabarc/go-catbox"
)

const DEFAULT_LABEL_TEXT = "No file selected"

type FileNameLabel struct {
	AbsPath string
	Label   *widget.Label
}

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
	CopyTextButton     *widget.Button
	CancelUploadButton *widget.Button
	FileName           *FileNameLabel
}

func NewUploadFileSection(app fyne.App, window fyne.Window, con *fyne.Container, absPath string) fyne.Widget {

	fileNameLabel := NewFileNameLabel(absPath)

	uploadBtn := widget.NewButtonWithIcon("", theme.UploadIcon(), nil)
	cancelBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), nil)
	openFileBtn := widget.NewButtonWithIcon("", theme.FolderNewIcon(), nil)
	copyTextBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if utils.PreferencesEnabled() {

			app.SendNotification(fyne.NewNotification("Copy", fmt.Sprintf("Copied: %s successfully", fileNameLabel.Label.Text)))
		}
		window.Clipboard().SetContent(fileNameLabel.Label.Text)
	})

	fileUploadWidget := NewFileUploadWidget(con, uploadBtn, cancelBtn, openFileBtn, copyTextBtn, fileNameLabel)
	return fileUploadWidget
}

func uploadToCatbox(fileName *FileNameLabel) {
	if fileName.Label.Text == DEFAULT_LABEL_TEXT {
		fmt.Println("no file selected, not uploading")
		return
	}
	file := fileName.Label.Text
	filePath := fileName.AbsPath
	fmt.Printf("uploading: %s\n", filePath)

	catboxClient := catbox.New(nil)
	catboxClient.Client.Timeout = time.Duration(fyne.CurrentApp().Preferences().IntWithFallback(utils.TIMEOUT_DURATION_MINUTES, utils.DEFAULT_FALLBACK_TIMEOUT_MINUTES)) * time.Minute
	if fyne.CurrentApp().Preferences().String(utils.CATBOX_USERHASH) != "" {
		catboxClient.Userhash = fyne.CurrentApp().Preferences().String(utils.CATBOX_USERHASH)
	}
	log.Printf("using userhash: %s", fyne.CurrentApp().Preferences().String(utils.CATBOX_USERHASH))
	if url, err := catboxClient.Upload(filePath); err != nil {
		fmt.Printf("\ncatbox: %v\n", err)
		fyne.CurrentApp().SendNotification(fyne.NewNotification("Error Uploading", err.Error()))
		return
	} else {
		embedLink := fmt.Sprintf("https://embeds.video/%s", url)
		fmt.Printf("\nurl: %s\npath: %s", embedLink, filePath)
		fileName.Label.SetText(embedLink)
		if utils.PreferencesEnabled() {
			title := fmt.Sprintf("%s Upload Finished", file)
			fyne.CurrentApp().SendNotification(fyne.NewNotification(title, ""))
		}
	}
}

func NewFileUploadWidget(container *fyne.Container, startUploadButton, cancelUploadButton, openFileButton, copyTextButton *widget.Button, fileName *FileNameLabel) *FileUploadWidget {

	item := &FileUploadWidget{
		OpenFileButton:     openFileButton,
		StartUploadButton:  startUploadButton,
		CopyTextButton:     copyTextButton,
		CancelUploadButton: cancelUploadButton,
		FileName:           fileName,
	}

	copyTextButton.Disable()
	cancelUploadButton.OnTapped = func() {

		fileName.Label.SetText(DEFAULT_LABEL_TEXT)
		if len(container.Objects) > 1 {
			fmt.Printf("length of container objects: %v\n", len(container.Objects))
			objPlace := len(container.Objects) - 1
			container.Remove(container.Objects[objPlace])
			return
		}
		if !copyTextButton.Disabled() {
			copyTextButton.Disable()
		}
		fmt.Println("file name cleared, nothing else happened")
	}

	startUploadButton.OnTapped = func() {
		go uploadToCatbox(fileName)
		if copyTextButton.Disabled() {
			copyTextButton.Enable()
		}
	}

	openFileButton.OnTapped = func() {
		fmt.Println("set open file tapped func")
		filename, err := dialog.File().Filter("Video Files", utils.VIDEO_FILE_EXTENSIOSN...).Load()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fileName.Label.SetText(filename)

		nextUploadBtn := widget.NewButtonWithIcon(startUploadButton.Text, startUploadButton.Icon, startUploadButton.OnTapped)
		nextCancelBtn := widget.NewButtonWithIcon(cancelUploadButton.Text, cancelUploadButton.Icon, cancelUploadButton.OnTapped)
		nextOpenFileBtn := widget.NewButtonWithIcon(openFileButton.Text, openFileButton.Icon, openFileButton.OnTapped)
		nextCopyTextBtn := widget.NewButtonWithIcon(copyTextButton.Text, copyTextButton.Icon, copyTextButton.OnTapped)

		container.Add(NewFileUploadWidget(container, nextUploadBtn, nextCancelBtn, nextOpenFileBtn, nextCopyTextBtn, NewFileNameLabel(DEFAULT_LABEL_TEXT)))
	}

	item.ExtendBaseWidget(item)
	return item
}

func (item *FileUploadWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.CancelUploadButton /* item.StartUploadButton, */, item.OpenFileButton, item.CopyTextButton, item.FileName.Label)
	return widget.NewSimpleRenderer(c)
}
