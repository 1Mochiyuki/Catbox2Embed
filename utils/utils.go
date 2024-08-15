package utils

import "fyne.io/fyne/v2"

const NOTIFICATIONS_ENABLED string = "notifications_enabled"

var VIDEO_FILE_EXTENSIOSN []string = []string{
	"mp4",
	"mov",
	"avi",
	"wmv",
	"flv",
	"webm",
	"gif",
	"m4v",
	"3gp",
	"mpeg",
	"mkv",
}

func PreferencesEnabled() bool {
	return fyne.CurrentApp().Preferences().Bool(NOTIFICATIONS_ENABLED)
}

