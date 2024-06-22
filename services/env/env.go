package env

import (
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/joho/godotenv"
)

type config struct {
  SESSION_SECRET string
  OAUTH2_GOOGLE_CLIENT_ID string
  OAUTH2_GOOGLE_CLIENT_SECRET string
}

var lock = &sync.Mutex{};
var configData *config;

func GetConfig() *config {
  if configData == nil {
    lock.Lock()
    defer lock.Unlock();

    if configData == nil {
      envFile, pathSet := os.LookupEnv("ENV_FILE");

      if !pathSet {
        envFile = ".env";
      }

      if err := godotenv.Load(envFile); err != nil {
        log.Fatal("Failed to load env file!");
      }
      
      configData = &config{};
      configStruct := reflect.ValueOf(configData).Elem();
      types := configStruct.Type();
      for i := 0; i < configStruct.NumField(); i++ {
        configStruct.Field(i).SetString(getEnvOrFail(types.Field(i).Name))
      }
    }
  }
  return configData;
}

func getEnvOrFail(key string) string {
  val, exists := os.LookupEnv(key);
  if !exists {
    log.Fatal(key + " not set!");
  }
  return val;
}
