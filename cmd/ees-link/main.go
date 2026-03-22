package main

import (
	"github.com/ksvaza/ees-link/envreader"
	"github.com/ksvaza/ees-link/logeris"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithError(errors.Errorf("%v", r)).Error("Panic")
		}
	}()

	err := envreader.SetupEnvreader(".env")
	if err != nil {
		logrus.WithError(err).Error("Failed to setup environment reader")
		return
	}

	logfile := envreader.GetEnvString("LOGFILE")
	clearlog := envreader.GetEnvBool("CLEARLOG")
	f, err := logeris.SetupLogger(logfile, clearlog)
	if err != nil {
		logrus.WithError(err).Error("Failed to setup logger")
		return
	}
	defer f.Close()

	logrus.Info("\nSveika, pasaule!\n")

	for i := 0; i < 5; i++ {
		logrus.Tracef("Trace %d", i)
		logrus.Debugf("Debug %d", i)
		logrus.Infof("Info %d", i)
		logrus.Warnf("Warn %d", i)
		logrus.Errorf("Error %d", i)
	}

	for {

	}
}
