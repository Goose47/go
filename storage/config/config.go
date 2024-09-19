package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var AppConfig struct {
	Mode string
	Host string
	Port string
}

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	AppConfig.Mode = checkAndRetrieve("APP_MODE")
	AppConfig.Host = checkAndRetrieve("APP_HOST")
	AppConfig.Port = checkAndRetrieve("APP_PORT")
}

func checkAndRetrieve(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("%s is not present in env.", key)
	}
	return val
}
