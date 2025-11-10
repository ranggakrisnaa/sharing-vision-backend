package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	// Default to info level
	Log.SetLevel(logrus.InfoLevel)
	// Output to stdout
	Log.SetOutput(os.Stdout)
	// Use JSON formatter for better structure in containers
	Log.SetFormatter(&logrus.JSONFormatter{})
}
