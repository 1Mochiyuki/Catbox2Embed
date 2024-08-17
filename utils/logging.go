package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logFile *os.File

func CreateLogFile() {
	CreateAppHome()
	logFileName := fmt.Sprintf("%s/info.log", AppHome())
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		os.Create(logFileName)
		fmt.Println(err)
	}
	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)
}

func CloseLogFile() {
	log.Println("Closing\n------------------------------------------------------")
	logFile.Close()
}