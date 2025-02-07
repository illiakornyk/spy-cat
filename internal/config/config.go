package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Env         string `yaml:"env" env-default:"development"`
    StoragePath string `yaml:"storage_path" env-required:"true"`
    HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
    Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
    Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
    IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        log.Fatal("CONFIG_PATH environment variable is not set")
    }

    if _, err := os.Stat(configPath); err != nil {
        log.Fatalf("error opening config file: %s", err)
    }

    var cfg Config

    err := readYAMLConfig(configPath, &cfg)
    if err != nil {
        log.Fatalf("error reading config file: %s", err)
    }

    return &cfg
}

func readYAMLConfig(path string, cfg *Config) error {
	fileContent, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    err = yaml.Unmarshal(fileContent, cfg)
    if err != nil {
        return err
    }

    return nil
}
