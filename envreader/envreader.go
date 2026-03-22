package envreader

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func viperGetString(name string) string {
	value, ok := viper.Get(name).(string)
	if ok {
		return value
	} else {
		return ""
	}
}
func viperGetInt(name string) int {
	s, ok := viper.Get(name).(string)
	if ok {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0
		} else {
			return n
		}
	} else {
		return 0
	}
}

func viperGetBool(name string) bool {
	s, ok := viper.Get(name).(string)
	if ok {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return false
		} else {
			return b
		}
	} else {
		return false
	}
}

func SetupEnvreader(envpath string) error {
	viper.SetConfigFile(envpath)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "Config read")
		}
	}
	return nil
}

func GetEnvString(name string) string {
	return viperGetString(name)
}

func GetEnvInt(name string) int {
	return viperGetInt(name)
}

func GetEnvBool(name string) bool {
	return viperGetBool(name)
}
