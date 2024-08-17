package utils

import (
	"fmt"
	"log"
	"os"
)

func AppHome() string {
	
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/Documents/Catbox2Embed/", userHome)
}

func CreateAppHome() {

	_, notFoundErr := os.Stat(AppHome())
	if notFoundErr != nil {
		os.Mkdir(AppHome(), 0755)
		log.Println("created app home")
	}
	log.Println("app home already created")
}