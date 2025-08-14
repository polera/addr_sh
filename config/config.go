package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tristanfisher/patchpanel"
	"reflect"
)

type Config struct {
	LogLevel string `default:"info"`

	ListenPort    string `default:":2000"`
	TLSListenPort string `default:":4443"`
	EnableTLS     bool   `default:"false"`
}

func ParseConfig(configPath string, configStruct Config) (*Config, error) {
	viperConf := viper.New()
	patch := patchpanel.NewPatchPanel(patchpanel.TokenSeparator, patchpanel.KeyValueSeparator)

	// get defaults off of our struct using patchpanel
	confType := reflect.TypeOf(configStruct)
	for i := 0; i < confType.NumField(); i++ {
		fieldVal, err := patch.GetDefault(confType.Field(i).Name, confType, []string{})
		if err != nil {
			return &Config{}, err
		}
		viperConf.SetDefault(confType.Field(i).Name, fieldVal)
	}

	// check configuration file for more values
	if configPath != "" {
		viperConf.SetConfigFile(configPath)
		err := viperConf.ReadInConfig()
		if err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				return nil, fmt.Errorf("file not found: %s", err)
			}
			return &Config{}, err
		}
	}

	viperConf.AutomaticEnv()
	err := viperConf.Unmarshal(&configStruct)
	if err != nil {
		return &Config{}, err
	}

	return &configStruct, nil
}
