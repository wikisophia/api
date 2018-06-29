package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	Addr                    string   `mapstructure:"addr"`
	ReadHeaderTimeoutMillis int      `mapstructure:"read_header_timeout_millis"`
	CorsAllowedOrigins      []string `mapstructure:"cors_allowed_origins"`
}

func (cfg *Server) ReadHeaderTimeout() time.Duration {
	return time.Duration(cfg.ReadHeaderTimeoutMillis) * time.Millisecond
}

func (cfg *Server) logValues() {
	log.Printf("server.addr=%s", cfg.Addr)
	log.Printf("server.read_header_timeout_millis=%d", cfg.ReadHeaderTimeoutMillis)
	log.Printf("server.cors_allowed_origins=%#v", cfg.CorsAllowedOrigins)
}

func (cfg *Server) setDefaults(v *viper.Viper) {
	v.SetDefault("server.addr", "localhost:8001")
	v.SetDefault("server.read_header_timeout_millis", 5000)
	v.SetDefault("server.cors_allowed_origins", []string{"*"})
}

func (cfg *Server) validate() []error {
	return validatePositiveInt("server.read_timeout_millis", cfg.ReadHeaderTimeoutMillis)
}

func validatePositiveInt(key string, value int) []error {
	if value <= 0 {
		return []error{fmt.Errorf(key+" must be a positive number. Got %d", value)}
	}
	return nil
}
