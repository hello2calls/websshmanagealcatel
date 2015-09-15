package logger

import (
	"log"
	"os"
)

// Print write log file with message
func Print(message string, err error) {
	logFile, errFile := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if errFile != nil {
		panic(errFile)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if err != nil {
		log.Print(message, err)
	} else {
		log.Print(message)
	}
}
