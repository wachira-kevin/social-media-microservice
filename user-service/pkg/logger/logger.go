package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
}
