package utils

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
)

const NOTIFICATIONS_ENABLED string = "notifications_enabled"
const CATBOX_USERHASH = "catbox_userhash"
const DEFAULT_FALLBACK_TIMEOUT_MINUTES = 30
const FILE_SIZE_REQUIREMENT = 200
const TIMEOUT_DURATION_MINUTES = "timeout_minutes"

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

func IsVideoFile(file string) bool {
	for _, v := range VIDEO_FILE_EXTENSIOSN {
		return strings.Contains(file, v)
	}
	fmt.Println("hello")
	return false
}

func PreferencesEnabled() bool {
	return fyne.CurrentApp().Preferences().Bool(NOTIFICATIONS_ENABLED)
}
