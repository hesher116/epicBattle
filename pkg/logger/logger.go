package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Fields map[string]interface{}

func init() {
	// Налаштування форматування
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func Info(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Info(message)
	} else {
		log.Info(message)
	}
}

func Warn(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Warn(message)
	} else {
		log.Warn(message)
	}
}

func Error(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Error(message)
	} else {
		log.Error(message)
	}
}

func Debug(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Debug(message)
	} else {
		log.Debug(message)
	}
}

// SetLevel встановлює рівень логування
func SetLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	}
}