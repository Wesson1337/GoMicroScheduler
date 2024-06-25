package logger

import (
	"log"
	"os"
)

func ConfigureLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
}