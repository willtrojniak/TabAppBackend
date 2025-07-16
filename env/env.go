package env

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	"github.com/joho/godotenv"
)

type config struct {
	BASE_URI                    string
	UI_URI                      string
	OAUTH2_GOOGLE_CLIENT_ID     string
	OAUTH2_GOOGLE_CLIENT_SECRET string
	OAUTH2_SLACK_CLIENT_ID      string
	OAUTH2_SLACK_CLIENT_SECRET  string
	POSTGRES_USER               string
	POSTGRES_PASSWORD           string
	POSTGRES_HOST               string
	POSTGRES_PORT               string
	POSTGRES_DB                 string
	REDIS_ADDR                  string
	EMAIL_CLIENT_HOST           string
	EMAIL_CLIENT_PORT           string
	EMAIL_CLIENT_USER           string
	EMAIL_CLIENT_PASSWORD       string
	EMAIL_CLIENT_ENABLED        bool
}

var Envs = getConfig()

const (
	DEV  string = "dev"
	PROD        = "prod"
)

var EXT_ENVIRONMENT string = DEV

func getConfig() config {
	envDir := os.Getenv("ENV_DIR")

	if err := godotenv.Load(filepath.Join(envDir, "base.env")); err != nil {
		log.Fatal("Failed to load base env file!")
	}

	if err := godotenv.Load(filepath.Join(envDir, EXT_ENVIRONMENT+".env")); err == nil {
		slog.Info("Loaded environment", "EXT_ENVIRONMENT", EXT_ENVIRONMENT)
	}

	configData := config{}
	configStruct := reflect.ValueOf(&configData).Elem()
	types := configStruct.Type()

	for i := range configStruct.NumField() {
		key := types.Field(i).Name
		switch configStruct.Field(i).Type().Kind() {
		case reflect.String:
			configStruct.Field(i).SetString(getEnvStringOrFail(key))
		case reflect.Bool:
			configStruct.Field(i).SetBool(getEnvBoolOrFail(key))
		default:
			log.Fatal(fmt.Sprintf("Unknown type for key '%v'", key))
		}
	}

	return configData
}

func getEnvStringOrFail(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal(key + " not set!")
	}
	return val
}

func getEnvBoolOrFail(key string) bool {
	str := getEnvStringOrFail(key)
	switch str {
	case "true":
		return true
	case "false":
		return false
	default:
		log.Fatal(fmt.Sprintf("%v not valid bool value ('%v')", key, str))
	}
	return false
}
