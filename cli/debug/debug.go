package debug

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Enable() {
	os.Setenv( "DEBUG", "1")
	logrus.SetLevel(logrus.DebugLevel)
}

func Disable() {
	os.Setenv("DEBUG", "")
	logrus.SetLevel(logrus.ErrorLevel)
}

func isEnabled() bool {
	return os.Getenv("DEBUG") != ""
}