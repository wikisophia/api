package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	ExternalURL        string   `mapstructure:"external_url"`
	Addr               string   `mapstructure:"addr"`
	ReadTimeoutMillis  int      `mapstructure:"read_timeout_millis"`
	WriteTimeoutMillis int      `mapstructure:"write_timeout_millis"`
	CorsAllowedOrigins []string `mapstructure:"cors_allowed_origins"`
}

func (cfg *Server) ReadTimeout() time.Duration {
	return time.Duration(cfg.ReadTimeoutMillis) * time.Millisecond
}

func (cfg *Server) WriteTimeout() time.Duration {
	return time.Duration(cfg.WriteTimeoutMillis) * time.Millisecond
}

func (cfg *Server) logValues() {
	log.Printf("server.addr=%s", cfg.Addr)
	log.Printf("server.external_url=%s", cfg.ExternalURL)
	log.Printf("server.read_timeout_millis=%d", cfg.ReadTimeoutMillis)
	log.Printf("server.write_timeout_millis=%d", cfg.WriteTimeoutMillis)
	log.Printf("server.cors_allowed_origins=%#v", cfg.CorsAllowedOrigins)
}

func (cfg *Server) setDefaults(v *viper.Viper) {
	log.Printf("server.addr=%s", "localhost:8001")
	v.SetDefault("server.external_url", "http://localhost:8001")
	v.SetDefault("server.read_timeout_millis", 5000)
	v.SetDefault("server.write_timeout_millis", 5000)
	v.SetDefault("server.cors_allowed_origins", []string{"*"})
}

func (cfg *Server) validate() []error {
	errs := validatePositiveInt("server.read_timeout_millis", cfg.ReadTimeoutMillis)
	errs = append(errs, validatePositiveInt("server.write_timeout_millis", cfg.WriteTimeoutMillis)...)
	return errs
}

func validatePositiveInt(key string, value int) []error {
	if value <= 0 {
		return []error{fmt.Errorf(key+" must be a positive number. Got %d", value)}
	}
	return nil
}
