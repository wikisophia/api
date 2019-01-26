package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// MustParse wraps Parse, but prints the errors and exits rather than returning them.
func MustParse() Configuration {
	cfg, errs := Parse()
	if len(errs) > 0 {
		log.Println("Invalid app config. Check your environment variables:")
		for _, err := range errs {
			log.Printf("  %v\n", err)
		}
		os.Exit(1)
	}
	return cfg
}

// Parse loads the app config using Viper.
// It logs all the values before returning, and panics on validation errors.
func Parse() (Configuration, []error) {
	// Config defaults go here
	cfg := Configuration{
		Server: &Server{
			Addr: "localhost:8001",
			ReadHeaderTimeoutMillis: 5000,
			CorsAllowedOrigins:      []string{"*"},
		},
		Storage: &Storage{
			Type: StorageTypeMemory,
			Postgres: &Postgres{
				Database: "wikisophia",
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "",
			},
		},
	}
	log.SetOutput(os.Stdout)
	errs := loadEnvironment(reflect.ValueOf(&cfg), "WKSPH_ARGS", "Configuration", nil)
	log.SetOutput(os.Stderr)

	errs = requirePositive(cfg.Server.ReadHeaderTimeoutMillis, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", errs)
	errs = requirePositive(cfg.Storage.Postgres.Port, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", errs)
	errs = requireValidStorageType(cfg.Storage.Type, "WKSPH_ARGS_STORAGE_TYPE", errs)
	return cfg, errs
}

func loadEnvironment(theValue reflect.Value, environmentVarSoFar string, pathSoFar string, errs []error) []error {
	theType := theValue.Type().Elem()

	for i := 0; i < theType.NumField(); i++ {
		thisField := theType.Field(i)
		thisFieldValue := theValue.Elem().Field(i)
		environment := environmentVarSoFar + "_" + thisField.Tag.Get("environment")
		path := pathSoFar + "." + thisField.Name
		switch thisField.Type.Kind() {
		case reflect.Ptr:
			errs = loadEnvironment(thisFieldValue, environment, path, errs)
		case reflect.Int, reflect.String:
			if value, isSet := os.LookupEnv(environment); isSet {
				errs = setSafely(environment, thisFieldValue, value, errs)
			}
			logIfLoggable(path, thisFieldValue)
		case reflect.Slice:
			if value, isSet := os.LookupEnv(environment); isSet {
				errs = setSliceSafely(thisField, thisFieldValue, value, errs)
			}
			logIfLoggable(path, thisFieldValue)
		default:
			panic("config.loadEnvironment() hasn't yet implemented parsing for type " + thisField.Type.String())
		}
	}
	return errs
}

func logIfLoggable(path string, value reflect.Value) {
	if strings.Contains(path, "Password") {
		log.Printf("%s: <redacted>", path)
	} else {
		log.Printf("%s: %#v", path, value)
	}
}

func setSafely(env string, toSet reflect.Value, value string, errs []error) []error {
	switch toSet.Kind() {
	case reflect.Int:
		errs = parseAndSetInt(env, toSet, value, errs)
	case reflect.String:
		toSet.SetString(value)
	default:
		panic(fmt.Sprintf("setSafely() is not yet implemented for type %v", toSet.Kind()))
	}

	return errs
}

func setSliceSafely(field reflect.StructField, toSet reflect.Value, value string, errs []error) []error {
	switch field.Type.Elem().Kind() {
	case reflect.String:
		toSet.Set(reflect.ValueOf(parseCommaSeparated(value)))
	default:
		panic(fmt.Sprintf("setSliceSafely() is not yet implement for slices of type %v", toSet.Kind()))
	}
	return errs
}

func parseAndSetInt(env string, toSet reflect.Value, value string, errs []error) []error {
	parsed, err := parseInt(value)
	if err != nil {
		return append(errs, fmt.Errorf("%s must be an int. Got \"%s\"", env, value))
	}
	if err == nil {
		toSet.SetInt(parsed)
		return errs
	}
	return append(errs, fmt.Errorf("%s: %v", env, err))
}

func parseInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func parseCommaSeparated(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

func requirePositive(value int, prefix string, errs []error) []error {
	if value <= 0 {
		return append(errs, fmt.Errorf("%s: must be positive. Got %d", prefix, value))
	}
	return errs
}

func requireValidStorageType(value StorageType, prefix string, errs []error) []error {
	allowedTypes := storageTypes()
	for _, storageType := range storageTypes() {
		if storageType == value {
			return errs
		}
	}
	return append(errs, fmt.Errorf("%s: must be one of %v. Got %s", prefix, allowedTypes, value))
}
