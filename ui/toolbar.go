package ui

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/1Mochiyuki/Catbox2Embed/utils"
)

func CreateToolBar(mainContainer *fyne.Container, a fyne.App, window fyne.Window, copyAllToolbarAction *widget.ToolbarAction) *widget.Toolbar {
	toolbar := widget.NewToolbar()

	notificationsOffStr := `<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M20 18.69L7.84 6.14 5.27 3.49 4 4.76l2.8 2.8v.01c-.52.99-.8 2.16-.8 3.42v5l-2 2v1h13.73l2 2L21 19.72l-1-1.03zM12 22c1.11 0 2-.89 2-2h-4c0 1.11.89 2 2 2zm6-7.32V11c0-3.08-1.64-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68c-.15.03-.29.08-.42.12-.1.03-.2.07-.3.11h-.01c-.01 0-.01 0-.02.01-.23.09-.46.2-.68.31 0 0-.01 0-.01.01L18 14.68z"/></svg>`
	notificationsOnStr := `<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.89 2 2 2zm6-6v-5c0-3.07-1.64-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.63 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/></svg>`
	baseNotificationOffResource := fyne.NewStaticResource("notifications_off", []byte(notificationsOffStr))
	baseNotificationOnResource := fyne.NewStaticResource("notifications_on", []byte(notificationsOnStr))

	notificationOnIcon := theme.NewThemedResource(baseNotificationOnResource)
	notificationOffIcon := theme.NewThemedResource(baseNotificationOffResource)

	var icon fyne.Resource

	if utils.PreferencesEnabled() {
		icon = notificationOnIcon
		a.Preferences().SetBool(utils.NOTIFICATIONS_ENABLED, true)
	} else {
		icon = notificationOffIcon
		a.Preferences().SetBool(utils.NOTIFICATIONS_ENABLED, false)
	}

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

	copyAllToolbarAction.Disable()

	helpToolbarAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		settingsWindow := a.NewWindow("Settings")
		settingsWindow.Resize(fyne.NewSize(400, 300))
		settingsWindow.SetFixedSize(true)

		userHash := a.Preferences().String(utils.CATBOX_USERHASH)
		userHashBinding := binding.BindString(&userHash)
		userHashEntry := widget.NewEntryWithData(userHashBinding)
		userHashEntry.SetPlaceHolder("Catbox Userhash")
		userHashEntry.OnSubmitted = func(text string) {
			a.Preferences().SetString(utils.CATBOX_USERHASH, text)
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
					a.Preferences().SetString(utils.CATBOX_USERHASH, text)
					fmt.Printf("userhash: %s\n", text)
					timer.Stop()
				}()
				return
			}
			a.Preferences().SetString(utils.CATBOX_USERHASH, text)
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

	createNewFileUploadSectionBtn := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		mainContainer.Add(NewUploadFileSection(a, window, mainContainer, DEFAULT_LABEL_TEXT))
	})

	toolbar.Append(createNewFileUploadSectionBtn)
	toolbar.Append(widget.NewToolbarSeparator())
	toolbar.Append(copyAllToolbarAction)
	toolbar.Append(widget.NewToolbarSpacer())
	toolbar.Append(notificationButton)
	toolbar.Append(helpToolbarAction)

	//content := container.NewBorder(toolbar, nil, nil, nil)
	//mainContainer.Add(content)
	uploadAllBtn := widget.NewButton("Upload All", func() {})
	uploadAllBtn.OnTapped = func() {
		for _, v := range mainContainer.Objects {

			if reflect.TypeOf(v) == reflect.TypeOf(&FileUploadWidget{}) {
				uploadWidget := v.(*FileUploadWidget)

				uploadWidget.StartUploadButton.OnTapped()
				if copyAllToolbarAction.Disabled() {
					copyAllToolbarAction.Enable()
				}
			}
		}
	}
	return toolbar
}
