package config

import (
    "os"
    "github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
    Logger = logrus.New()
    Logger.SetOutput(os.Stdout)
    Logger.SetLevel(logrus.InfoLevel)
    Logger.SetFormatter(&logrus.JSONFormatter{})
}
