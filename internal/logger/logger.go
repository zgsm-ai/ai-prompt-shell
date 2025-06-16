package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(level string) {
	log = logrus.New()
	log.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func GetLogger() *logrus.Logger {
	return log
}
