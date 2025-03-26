package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Auth        `yaml:"auth"`
}

type Auth struct {
	User     string `yaml:"user" env:"AUTH_USER"`
	Password string `yaml:"password" env:"AUTH_PASSWORD"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	fmt.Println(err)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// 3. Явная загрузка из переменных окружения (может перезаписать значения из YAML)
	if user := os.Getenv("AUTH_USER"); user != "" {
		cfg.Auth.User = user
		log.Println("Using AUTH_USER from environment")
	}

	if password := os.Getenv("AUTH_PASSWORD"); password != "" {
		cfg.Auth.Password = password
		log.Println("Using AUTH_PASSWORD from environment")
	}

	if cfg.Auth.User == "" || cfg.Auth.Password == "" {
		log.Fatal("Username and Password is empty")
	}

	return &cfg
}
