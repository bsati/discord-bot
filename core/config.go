package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents config variables needed in different parts of the app
type Config struct {
	DbConnectionString string
	BotToken           string
}

// LoadConfig loads a Config from the given filepath, if specified, else from
// environment variables.
func LoadConfig(filepath *string) Config {
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
	var connectionString string
	if dbName == "" {
		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass)
	} else {
		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	}
	botToken := os.Getenv("BOT_TOKEN")
	return &Config{
		DbConnectionString: connectionString,
		BotToken:           botToken,
	}
}
