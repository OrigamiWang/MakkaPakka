package logutil

import (
	"log"
)

func Info(format string, v ...any) {
	log.Printf(format, v...)
}
