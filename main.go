package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/1Mochiyuki/Catbox2Embed/fileupload"
	"github.com/1Mochiyuki/Catbox2Embed/utils"
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
		var linksCount int
		for _, v := range con.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&fileupload.FileUploadWidget{}) {
				uploadWidget := v.(*fileupload.FileUploadWidget)
				if uploadWidget.FileName.Label.Text == fileupload.DEFAULT_LABEL_TEXT {
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

func newUploadFileSection(app fyne.App, window fyne.Window, con *fyne.Container, absPath string) fyne.Widget {

	fileNameLabel := fileupload.NewFileNameLabel(absPath)

	uploadBtn := widget.NewButtonWithIcon("", theme.UploadIcon(), nil)
	cancelBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), nil)
	openFileBtn := widget.NewButtonWithIcon("", theme.FolderNewIcon(), nil)
	copyTextBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if utils.PreferencesEnabled() {

			app.SendNotification(fyne.NewNotification("Copy", fmt.Sprintf("Copied: %s successfully", fileNameLabel.Label.Text)))
		}
		window.Clipboard().SetContent(fileNameLabel.Label.Text)
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

const FILE_SIZE_REQUIREMENT = 200

func isVideoFile(file string) bool {
	for _, v := range utils.VIDEO_FILE_EXTENSIOSN {
		return strings.Contains(file, v)
	}
	fmt.Println("hello")
	return false
}

func main() {
	notificationsOffStr := `<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M20 18.69L7.84 6.14 5.27 3.49 4 4.76l2.8 2.8v.01c-.52.99-.8 2.16-.8 3.42v5l-2 2v1h13.73l2 2L21 19.72l-1-1.03zM12 22c1.11 0 2-.89 2-2h-4c0 1.11.89 2 2 2zm6-7.32V11c0-3.08-1.64-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68c-.15.03-.29.08-.42.12-.1.03-.2.07-.3.11h-.01c-.01 0-.01 0-.02.01-.23.09-.46.2-.68.31 0 0-.01 0-.01.01L18 14.68z"/></svg>`
	notificationsOnStr := `<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.89 2 2 2zm6-6v-5c0-3.07-1.64-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.63 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/></svg>`
	baseNotificationOffResource := fyne.NewStaticResource("notifications_off", []byte(notificationsOffStr))
	baseNotificationOnResource := fyne.NewStaticResource("notifications_on", []byte(notificationsOnStr))

	notificationOnIcon := theme.NewThemedResource(baseNotificationOnResource)
	notificationOffIcon := theme.NewThemedResource(baseNotificationOffResource)
	var icon fyne.Resource

	a := app.NewWithID("Catbox2Embed")

	fmt.Printf("notifications: %v\n", utils.PreferencesEnabled())
	fmt.Printf("catbox userhash: %s\n", a.Preferences().String(fileupload.CATBOX_USERHASH))

	fmt.Println("setting pref")
	if utils.PreferencesEnabled() {
		icon = notificationOnIcon
		a.Preferences().SetBool(utils.NOTIFICATIONS_ENABLED, true)
	} else {
		icon = notificationOffIcon
		a.Preferences().SetBool(utils.NOTIFICATIONS_ENABLED, false)
	}

	window := a.NewWindow("Catbox2Embed")
	window.Resize(fyne.NewSize(600, 500))

	notificationButton := widget.NewToolbarAction(icon, func() {})
	notificationButton.OnActivated = func() {

		if utils.PreferencesEnabled() {
			fmt.Println("notifications off")
			a.Preferences().SetBool(utils.NOTIFICATIONS_ENABLED, false)
			notificationButton.SetIcon(notificationOffIcon)
			return
		} else {
			fmt.Println("notifications on")
			a.Preferences().SetBool(utils.NOTIFICATIONS_ENABLED, true)
			notificationButton.SetIcon(notificationOnIcon)
			return
		}
	}

	toolbar := widget.NewToolbar()

	mainContainer := container.NewVBox()
	createNewFileUploadSectionBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		mainContainer.Add(newUploadFileSection(a, window, mainContainer, fileupload.DEFAULT_LABEL_TEXT))
	})
	copyAllToolbarAction := widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		var links string
		for _, v := range mainContainer.Objects {

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
	copyAllToolbarAction.Disable()

	helpToolbarAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		settingsWindow := a.NewWindow("Settings")
		settingsWindow.Resize(fyne.NewSize(400, 300))
		settingsWindow.SetFixedSize(true)

		userHash := a.Preferences().String(fileupload.CATBOX_USERHASH)
		userHashBinding := binding.BindString(&userHash)
		userHashEntry := widget.NewEntryWithData(userHashBinding)
		userHashEntry.SetPlaceHolder("Catbox Userhash")
		userHashEntry.OnSubmitted = func(text string) {
			a.Preferences().SetString(fileupload.CATBOX_USERHASH, text)
		}
		userHashEntry.OnChanged = func(text string) {
			textLen := len(text)
			if textLen > 30 {
				fmt.Println("userhash too long")
				userHashEntry.SetValidationError(errors.New("userhash too long"))
				return
			}
			if textLen >= 1 {
				timer := time.NewTimer(time.Millisecond * 150)
				go func() {
					<-timer.C
					a.Preferences().SetString(fileupload.CATBOX_USERHASH, text)
					fmt.Printf("userhash: %s\n", text)
					timer.Stop()
				}()
				return
			}
			a.Preferences().SetString(fileupload.CATBOX_USERHASH, text)
			fmt.Println("userhash blank")
		}

		con := container.NewVBox(userHashEntry)
		shortcutsLabel := widget.NewLabel("Shortcuts:\nCopy all links: Ctrl + E")
		shortcutsLabel.Alignment = fyne.TextAlignCenter
		shortcutsLabel.TextStyle.Bold = true
		con.Add(shortcutsLabel)

		settingsWindow.SetContent(con)
		settingsWindow.Show()
	})
	toolbar.Append(createNewFileUploadSectionBtn)
	toolbar.Append(widget.NewToolbarSeparator())
	toolbar.Append(copyAllToolbarAction)
	toolbar.Append(widget.NewToolbarSpacer())
	toolbar.Append(notificationButton)
	toolbar.Append(helpToolbarAction)

	content := container.NewBorder(toolbar, nil, nil, nil)
	mainContainer.Add(content)
	uploadAllBtn := widget.NewButton("Upload All", func() {})
	uploadAllBtn.OnTapped = func() {
		for _, v := range mainContainer.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&fileupload.FileUploadWidget{}) {
				uploadWidget := v.(*fileupload.FileUploadWidget)

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
			if sizeInMib > int64(FILE_SIZE_REQUIREMENT) {
				dialog.ShowError(fmt.Errorf("size must be %v MiB or lower. file was: %v MiB", FILE_SIZE_REQUIREMENT, sizeInMib), window)
				continue
			}
			if !isVideoFile(fileInfo.Name()) {
				dialog.ShowError(errors.New("file is not a video.\nif this is an error, please contact the developer"), window)
				continue
			}

			uploadWidget := newUploadFileSection(a, window, mainContainer, v.Path())
			mainContainer.Add(uploadWidget)
		}
		if len(mainContainer.Objects) > 1 {
			addCopyAllLinksShortcut(window, mainContainer)
			scrollContainer := container.NewVScroll(mainContainer)
			window.SetContent(scrollContainer)
		}

	})

	if len(mainContainer.Objects) >= 0 {

		instructions := NewInstructions(fmt.Sprintf("Click or drag files to begin. Must be %v MiB or lower\nPress Ctrl + E while window is focused to copy all links", FILE_SIZE_REQUIREMENT), func() {
			fileUploadSection := newUploadFileSection(a, window, mainContainer, fileupload.DEFAULT_LABEL_TEXT)
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
