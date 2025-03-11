package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DBHost     string `yaml:"db_host"`
	DBPort     string `yaml:"db_port"`
	DBUsername string `yaml:"db_username"`
	DBName     string `yaml:"db_name"`
	DBPassword string `yaml:"db_password"`
	Port       string `yaml:"port"`
	JWTSecret  string `yaml:"jwt_secret"`
}

var (
	conf     *Config
	confOnce sync.Once
)

func GetConf() *Config {
	confOnce.Do(func() {
		configPath := "conf.yaml"

		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			configPath = envPath
		}

		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		conf = &Config{}
		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			log.Fatalf("Error parsing config file: %v", err)
		}

		if conf.DBHost == "" || conf.DBPort == "" || conf.DBName == "" {
			log.Fatalf("Missing required database configuration")
		}

		if conf.JWTSecret == "" {
			log.Fatalf("Missing required JWT secret configuration")
		}

		log.Printf("Configuration loaded successfully from %s", configPath)
	})

	return conf
}
