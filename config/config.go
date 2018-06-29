package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Server  Server  `mapstructure:"server"`
	Storage Storage `mapstructure:"storage"`
}

func (cfg *Configuration) logValues() {
	cfg.Server.logValues()
	cfg.Storage.logValues()
}

func (cfg *Configuration) setDefaults(v *viper.Viper) {
	cfg.Server.setDefaults(v)
	cfg.Storage.setDefaults(v)
}

func (cfg *Configuration) validate() []error {
	errs := cfg.Server.validate()
	errs = append(errs, cfg.Storage.validate()...)
	return errs
}
