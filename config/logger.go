// config/logger.go

package config

import (
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(
		os.Stdout,
		"[BHERI] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}
