package log

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/suyashkumar/conduit/config"
)

// Configure logging for this project
func Configure() {
	f := config.Get(config.LogFile)
	h := lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  f,
		logrus.WarnLevel:  f,
		logrus.ErrorLevel: f,
		logrus.DebugLevel: f,
		logrus.FatalLevel: f,
		logrus.PanicLevel: f,
	}, &logrus.JSONFormatter{})

	logrus.AddHook(h)
}
