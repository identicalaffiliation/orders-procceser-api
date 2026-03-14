package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServiceConfig struct {
	PostgresConfig postgresConfig
	LoggerConfig   loggerConfig `yaml:"logger"`
	RabbitMQConfig rabbitConfig
	ServerConfig   serverConfig  `yaml:"api"`
	Timeout        time.Duration `yaml:"timeout"`
}

type postgresConfig struct {
	Port     string `env:"POSTGRES_PORT"`
	Host     string `env:"POSTGRES_HOST"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBname   string `env:"POSTGRES_DB"`
	SSLmode  string `env:"POSTGRES_SSL_MODE"`
}

type loggerConfig struct {
	Level string `yaml:"level"`
}

type rabbitConfig struct {
	URI      string `env:"RABBITMQ_URI"`
	Exchange string `env:"RABBITMQ_EXCHANGE"`
}

type serverConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

func MustLoadConfig(configPath string) *ServiceConfig {
	cfg := new(ServiceConfig)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(err)
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic(err)
	}

	return cfg
}
