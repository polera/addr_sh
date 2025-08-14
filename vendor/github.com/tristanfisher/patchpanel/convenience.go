package patchpanel

import (
	"flag"
	"os"
)

const ENV_CONFIG_FILE = "CONFIG_FILE"
const FLAG_CONFIG_FILE = "config_file"

func GetFileEnvOrPath(envKey string, flagKey string) string {
	valueFilePath := os.Getenv(envKey)
	valueFilePathFlag := flag.String(flagKey, "", "path to target file of values")
	flag.Parse()
	// command line value takes precedence over environment variable
	if valueFilePathFlag != nil && len(*valueFilePathFlag) > 0 {
		valueFilePath = *valueFilePathFlag
	}
	return valueFilePath
}
