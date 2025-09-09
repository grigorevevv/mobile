package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(level)
	return logger
}
