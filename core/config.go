package core

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbConnectionString string
	BotToken           string
}

func loadConfig(filepath *string) Config {
	if filepath == nil {
		return *loadConfigFromEnv()
	}
	cfg, err := loadConfigFromFile(*filepath)
	if err != nil {
		panic(err)
	}
	return *cfg
}

func loadConfigFromFile(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	cfg := &Config{}
	err = decoder.Decode(cfg)
	return cfg, err
}

func loadConfigFromEnv() *Config {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	botToken := os.Getenv("BOT_TOKEN")
	return &Config{
		DbConnectionString: connectionString,
		BotToken:           botToken,
	}
}
