package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const envPrefix = "WKSPH_ARGS"

// ParseConfig loads the app config using Viper.
// It logs all the values before returning, and panics on validation errors.
func ParseConfig() Configuration {
	return ParseConfigFromPath(".")
}

// ParseConfigFromPath helps unit tests so they can use the same config as prod.
func ParseConfigFromPath(path string) Configuration {
	var cfg Configuration
	v := viper.New()

	cfg.setDefaults(v)
	respectConfigFile(v, path)
	respectEnvironmentVariables(v)

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
	cfg.logValues()
	panicOnErrors(cfg.validate())
	return cfg
}

func respectConfigFile(v *viper.Viper, path string) {
	v.SetConfigName("config")
	v.AddConfigPath(path)
	if err := v.ReadInConfig(); err != nil {
		// If the config file doesn't exist, that should be fine. It'll just use
		// the defaults & environment variables.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Failed to read the config: %v", err)
		}
	}
}

func respectEnvironmentVariables(v *viper.Viper) {
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

func panicOnErrors(errs []error) {
	if len(errs) > 0 {
		log.Println("Error: invalid app config. Check your config.yaml or " + envPrefix + "_* environment variables.")
		for _, err := range errs {
			log.Println("  " + err.Error())
		}
		os.Exit(1)
	}
}
