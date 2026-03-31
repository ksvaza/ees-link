package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/envreader"
	"github.com/ksvaza/ees-link/httpapi"
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

	// for i := 0; i < 5; i++ {
	// 	logrus.Tracef("Trace %d", i)
	// 	logrus.Debugf("Debug %d", i)
	// 	logrus.Infof("Info %d", i)
	// 	logrus.Warnf("Warn %d", i)
	// 	logrus.Errorf("Error %d", i)
	// }

	// sudi

	ctx := context.Background()

	slices := make([]db.RandomStruct, 1000)
	slices2 := make([]db.Ahh, 1000)

	for i := 0; i < 1000; i++ {
		slices[i] = db.RandomStruct{
			ID:         i + 1,
			InstanceID: (i % 10) + 1,
			Name:       fmt.Sprintf("Name%d", i+1),
			Value:      float64(i) * 1.1,
			Timestamp:  time.Now(),
		}
		slices2[i] = db.Ahh{
			ID:         i + 1,
			InstanceID: (i % 10) + 1,
			Name:       fmt.Sprintf("Name%d", i+1),
			Value:      float64(i) * 1.1,
			Timestamp:  time.Now(),
			Nuniga:     fmt.Sprintf("Nuniga%d", i+1),
		}
	}

	//strukti1 := db.ToAnySlice(slices)
	//strukti2 := db.ToAnySlice(slices2)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envreader.GetEnvString("POSTGRES_HOST"),
		envreader.GetEnvString("POSTGRES_PORT"),
		envreader.GetEnvString("POSTGRES_USER"),
		envreader.GetEnvString("POSTGRES_PASSWORD"),
		envreader.GetEnvString("POSTGRES_DB"),
	)

	db.Init(dsn)

	//db.BatchInsertIntoTable(ctx, "kkas", strukti1)
	// db.BatchInsertIntoTable(ctx, "kkas2", strukti2)

	results, err := db.ReadFromTableWhere(ctx, "kkas", db.RandomStruct{}, "instance_i_d", []any{1, 2, 3})
	if err != nil {
		logrus.WithError(err).Error("Failed to read from table")
		return
	}

	for i, r := range results {
		v := reflect.ValueOf(r)
		t := reflect.TypeOf(r)

		fmt.Printf("--- Row %d ---\n", i+1)
		for j := 0; j < t.NumField(); j++ {
			fmt.Printf("%s: %v\n", t.Field(j).Name, v.Field(j).Interface())
		}
	}

	fmt.Println("success")

	logrus.Info("\nSveika, http aplikācija!\n")

	err = httpapi.SetupHTTPAPI()
	if err != nil {
		logrus.WithError(errors.Wrap(err, "HTTP")).Error("Error")
	}

	for {

	}
}
