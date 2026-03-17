package main

import (
	"fmt"

	"github.com/ksvaza/ees-link/logeris"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithError(errors.New(fmt.Sprintf("%v", r))).Error("Panic")
		}
	}()

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "02.01.2006 15:04:05.000",
		FullTimestamp:   true,
		//		DisableColors:   true,
		DisableQuote: true,
	})
	logrus.SetOutput(&logeris.LogWriter{})

	logrus.AddHook(&logeris.StacktraceHook{})

	logrus.SetLevel(logrus.InfoLevel)

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
