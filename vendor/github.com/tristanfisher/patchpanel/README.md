# patchpanel


`patchpanel` aims to make it easy to load up a struct with data at runtime, such as loading up configuration data. 

no external dependencies are used.

### functionality

- getting values via [struct tags](https://go.dev/ref/spec#Tag)
- type coercions / deserializers

### example usage

example using viper in conjunction with patchpanel to load configuration onto a struct

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tristanfisher/patchpanel"
	"os"
	"reflect"
	"time"
)

type Config struct {
	Runtime time.Duration `default:"5s"`
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

func main() {
	// grab a values/configuration file path from our environment using patchpanel
	valuesFile := patchpanel.GetFileEnvOrPath(patchpanel.ENV_CONFIG_FILE, patchpanel.FLAG_CONFIG_FILE)

	conf, err := ParseConfig(valuesFile, Config{})
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("error parsing configuration: %s\n", err.Error()))
		os.Exit(1)
	}

	_, _  = os.Stdout.WriteString(fmt.Sprintf(":)\n"))
	ctx, cancel := context.WithTimeout(context.Background(), conf.Runtime)
	defer cancel()
	select {
		case <-ctx.Done():
			fmt.Println("have a good day")
	}
}
```