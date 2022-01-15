package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"avito-tech-task/config"
)

func NewLogger(config *config.Config) (*logrus.Logger, func() error) {
	level, err := logrus.ParseLevel(config.LoggingLevel)
	if err != nil {
		logrus.Fatal("Could not parse logging level: %s", err)
	}

	logger := logrus.New()

	format := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second()) + ".log"

	file, err := os.OpenFile(format, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Could not open file %s: %s", format, err)
	}

	logger.Writer()
	logger.SetLevel(level)
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger, file.Close
}
