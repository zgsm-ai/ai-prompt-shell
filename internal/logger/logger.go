package logger

import (
	"github.com/zgsm-ai/ai-prompt-shell/internal/config"
	"os"

	"github.com/sirupsen/logrus"
)

func Init(cfg *config.LoggerConfig) {
	// Set log level
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = logrus.DebugLevel
	}
	logrus.SetLevel(lvl)

	// Set log format
	switch cfg.LogFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default: // text as default
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// Set log output
	switch cfg.LogOutput {
	case "file":
		file, err := os.OpenFile(cfg.LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(file)
		} else {
			logrus.SetOutput(os.Stdout)
			logrus.Warn("Failed to log to file, using default stdout")
		}
	default: // stdout as default
		logrus.SetOutput(os.Stdout)
	}
}
