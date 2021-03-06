package logger

import (
	"github.com/cloudfoundry-community/bosh-vault/config"
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func Initialize(bvConfig config.Configuration) {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	logLevel, err := logrus.ParseLevel(bvConfig.Log.Level)
	if err != nil {
		Log.Errorf("error parsing configured log level %s, defaulting to debug", bvConfig.Log.Level)
		logLevel = logrus.DebugLevel
	}
	Log.SetLevel(logLevel)
	Log.Out = os.Stdout
}
