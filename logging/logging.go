package logging

import (
	"log"
	"os"
)

var (
	InfoLog *log.Logger
)

func InitLogging() {
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(file, "MESSAGE: ", log.Lshortfile)

}
