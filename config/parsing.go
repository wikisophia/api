package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// MustParseConfig wraps ParseConfig, but prints the errors and exits rather than returning them.
func MustParseConfig() Configuration {
	cfg, errs := ParseConfig()
	if len(errs) > 0 {
		log.Println("Invalid app config. Check your environment variables:")
		for _, err := range errs {
			log.Printf("  %v\n", err)
		}
		os.Exit(1)
	}
	return cfg
}

// ParseConfig loads the app config using Viper.
// It logs all the values before returning, and panics on validation errors.
func ParseConfig() (Configuration, []error) {
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

	errs := loadEnvironment(reflect.ValueOf(&cfg), envPrefix, "Configuration", nil)

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
			loadEnvironment(thisFieldValue, environment, path, errs)
		case reflect.Int, reflect.String:
			errs = setSafely(thisFieldValue, os.Getenv(environment), errs)
			logIfLoggable(path, thisFieldValue)
		case reflect.Slice:
			errs = setSliceSafely(thisField, thisFieldValue, os.Getenv(environment), errs)
			logIfLoggable(path, thisFieldValue)
		default:
			panic("config.loadEnvironment() hasn't yet implemented parsing for type " + thisField.Type.String())
		}
	}
	return errs
}

func logIfLoggable(path string, value reflect.Value) {
	if strings.Contains(path, "Password") {
		log.Printf("%s: <redacted for security>", path)
	} else {
		log.Printf("%s: %#v", path, value)
	}
}

func setSafely(toSet reflect.Value, value string, errs []error) []error {
	if value == "" {
		return errs
	}

	switch toSet.Kind() {
	case reflect.Int:
		errs = parseAndSetInt(toSet, value, errs)
	case reflect.String:
		toSet.SetString(value)
	default:
		panic(fmt.Sprintf("setSafely() is not yet implemented for type %v", toSet.Kind()))
	}

	return errs
}

func setSliceSafely(field reflect.StructField, toSet reflect.Value, value string, errs []error) []error {
	if value == "" {
		return errs
	}

	switch field.Type.Elem().Kind() {
	case reflect.String:
		toSet.Set(reflect.ValueOf(parseSpaceSeparated(value)))
	default:
		panic(fmt.Sprintf("setSliceSafely() is not yet implement for slices of type %v", toSet.Kind()))
	}
	return errs
}

func parseAndSetInt(toSet reflect.Value, value string, errs []error) []error {
	parsed, err := parseInt(value)
	if err == nil {
		toSet.SetInt(parsed)
	}
	return append(errs, err)
}

func parseInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func parsePositiveInt(value string) (int, error) {
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil || parsed <= 0 {
		return -1, fmt.Errorf("%s was not a positive int", value)
	}
	return int(parsed), nil
}

func parseSpaceSeparated(value string) []string {
	return strings.Split(value, " ")
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
