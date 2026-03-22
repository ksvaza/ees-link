package main

import (
	"fmt"
	"os"

	"github.com/ksvaza/ees-link/logeris"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func setupLogger(logFile string) (f *os.File, err error) {

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "02.01.2006 15:04:05.000",
		FullTimestamp:   true,
		//		DisableColors:   true,
		DisableQuote: true,
	})

	f, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	// Output to stdout instead of the default stderr
	logrus.SetOutput(f)

	//logrus.SetOutput(&logeris.LogWriter{})

	logrus.AddHook(&logeris.StacktraceHook{})

	logrus.SetLevel(logrus.InfoLevel)

	return
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithError(errors.New(fmt.Sprintf("%v", r))).Error("Panic")
		}
	}()

	logfile := "./log.txt"
	f, err := setupLogger(logfile)
	if err != nil {
		fmt.Println("Failed to create logfile" + logfile)
		panic(err)
	}
	defer f.Close()

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
