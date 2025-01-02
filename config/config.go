package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host              string `yaml:"host"`
	Port              string `yaml:"port"`
	Basepath          string `yaml:"basepath"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	Online_path       string `yaml:"online_path"`
	AppPort           string `yaml:"app_port"`
	BasicAuthUsername string `yaml:"basic_auth_username"`
	BasicAuthPassword string `yaml:"basic_auth_password"`
}

func LoadConfig() *Config {
	config := &Config{
		Host:              getEnv("HOST", ""),
		Port:              getEnv("PORT", ""),
		Basepath:          getEnv("BASEPATH", ""),
		Username:          getEnv("USERNAME", ""),
		Password:          getEnv("PASSWORD", ""),
		AppPort:           getEnv("APP_PORT", "8080"),
		Online_path:       "panel/api/inbounds/onlines",
		BasicAuthUsername: getEnv("BASIC_AUTH_USERNAME", "admin"),
		BasicAuthPassword: getEnv("BASIC_AUTH_PASSWORD", "admin"),
	}

	missingVars := []string{}
	if config.Host == "" {
		missingVars = append(missingVars, "HOST")
	}
	if config.Port == "" {
		missingVars = append(missingVars, "PORT")
	}
	if config.Basepath == "" {
		missingVars = append(missingVars, "BASEPATH")
	}
	if config.Username == "" {
		missingVars = append(missingVars, "USERNAME")
	}
	if config.Password == "" {
		missingVars = append(missingVars, "PASSWORD")
	}
	if config.AppPort == "" {
		missingVars = append(missingVars, "APP_PORT")
	}
	if config.BasicAuthUsername == "" {
		missingVars = append(missingVars, "BASIC_AUTH_USERNAME")
	}
	if config.BasicAuthPassword == "" {
		missingVars = append(missingVars, "BASIC_AUTH_PASSWORD")
	}
	if len(missingVars) > 0 {
		log.Printf("missing required environment variables: %v", missingVars)
		return FileConfig()
	}
	log.Printf("config loaded from environment variables")
	return config
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func FileConfig() *Config {
	log.Println("loading config from config.yaml file")
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error opening config.yaml file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalf("error decoding config.yaml file: %v", err)
	}
	missingVars := []string{}
	if config.Host == "" {
		missingVars = append(missingVars, "HOST")
	}
	if config.Port == "" {
		missingVars = append(missingVars, "PORT")
	}
	if config.Basepath == "" {
		missingVars = append(missingVars, "BASEPATH")
	}
	if config.Username == "" {
		missingVars = append(missingVars, "USERNAME")
	}
	if config.Password == "" {
		missingVars = append(missingVars, "PASSWORD")
	}
	if config.AppPort == "" {
		missingVars = append(missingVars, "APP_PORT")
	}
	if config.BasicAuthUsername == "" {
		missingVars = append(missingVars, "BASIC_AUTH_USERNAME")
	}
	if config.BasicAuthPassword == "" {
		missingVars = append(missingVars, "BASIC_AUTH_PASSWORD")
	}
	if config.Online_path == "" {
		config.Online_path = "panel/api/inbounds/onlines"
	}
	if len(missingVars) > 0 {
		log.Fatalf("config.yaml missing following entries: %v", missingVars)
	}

	log.Println("config loaded from config.yaml file")
	return config
}
