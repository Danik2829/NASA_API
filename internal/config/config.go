package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB          PsqlConfig
	MinioConfig MinioConfig
}

type PsqlConfig struct {
	DBDriver   string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
		return nil, err
	}

	psqlcfg := &PsqlConfig{
		DBDriver:   os.Getenv("POSTGRES_DRIVER"),
		DBPort:     os.Getenv("POSTGRES_PORT"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBHost:     os.Getenv("POSTGRES_HOST"),
	}

	miniocfg := &MinioConfig{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Bucket:    os.Getenv("MINIO_BUCKET_NAME"),
	}

	cfg := &Config{
		DB:          *psqlcfg,
		MinioConfig: *miniocfg,
	}

	return cfg, nil
}
