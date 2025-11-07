package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type AppConfig struct {
	ApplicatonName string
	AppPort        string
}

type ApiConfig struct {
	JwtSignaturKey 		string
	JwtSigningMethod 	string
	AccessTokenLifeTime int
}

type Config struct {
	DBConfig
	AppConfig
	ApiConfig
}

func (cfg *Config) loadConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg.AppConfig = AppConfig{
		ApplicatonName: os.Getenv("APP_NAME"),
		AppPort:        os.Getenv("APP_PORT"),
	}

	cfg.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	cfg.ApiConfig = ApiConfig{
		JwtSignaturKey:     os.Getenv("JWT_SIGNATURE_KEY"),
		JwtSigningMethod:   os.Getenv("JWT_SIGNING_METHOD"),
		AccessTokenLifeTime: 24,
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.Database == "" || cfg.Username == "" || cfg.Password == "" || cfg.AppPort == "" || cfg.JwtSignaturKey == "" {
		fmt.Println("config .env is required")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	
	err := cfg.loadConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
