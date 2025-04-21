package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func InitLogger() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Logger.Fatal("Error opening log file:", err)
	}
	Logger.SetOutput(logFile)
	Logger.SetLevel(logrus.InfoLevel)

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
