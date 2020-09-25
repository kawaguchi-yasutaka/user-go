package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	UserName     string
	UserPassword string
	Address      string
	Port         string
	Table        string
}

func SetConfig() Config {
	env := os.Getenv("USER_GO_ENV")
	loadEnv(env)
	config := Config{}
	config.setDBConfig()
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
