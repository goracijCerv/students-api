package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address"`
}

type SMTPConfiguration struct {
	From     string `yaml:"from"`
	Password string `yaml:"password"`
	SmtpHost string `yaml:"smtpHost"`
	SmtpPort string `yaml:"smtpPort"`
}

// env-default:"production"
type Config struct {
	Env               string `yaml:"env" env:"ENV" env-required:"true" `
	StoragePath       string `yaml:"storage_path" env-required:"true"`
	HTTPServer        `yaml:"http_server"`
	SMTPConfiguration `yaml:"smtp"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set ")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file foest not exist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}
	return &cfg
}
