package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DB  DBConfig
	APP APPConfig
}

type DBConfig struct {
	UserName     string
	UserPassword string
	Address      string
	Port         string
	Table        string
}

type APPConfig struct {
	URL           string
	EmailAddress  string
	EmailPassword string
	EmailHost     string
}

func SetConfig() Config {
	env := os.Getenv("USER_GO_ENV")
	loadEnv(env)
	config := Config{}
	config.setDBConfig()
	config.setAppConfig()
	return config
}

func loadEnv(env string) {
	if env == "" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(fmt.Sprintf(".env.%s", env))
}

func (config *Config) setDBConfig() {
	config.DB.UserName = os.Getenv("DB_USER_NAME")
	config.DB.UserPassword = os.Getenv("DB_USER_PASSWORD")
	config.DB.Address = os.Getenv("DB_ADDRESS")
	config.DB.Port = os.Getenv("DB_PORT")
	config.DB.Table = os.Getenv("DB_TABLE")
}

func (config *Config) setAppConfig() {
	config.APP.URL = os.Getenv("APP_URL")
	config.APP.EmailAddress = os.Getenv("APP_EMAIL_ADDRESS")
	config.APP.EmailPassword = os.Getenv("APP_EMAIL_PASSWORD")
	config.APP.EmailHost = os.Getenv("APP_EMAIL_HOST")
}
