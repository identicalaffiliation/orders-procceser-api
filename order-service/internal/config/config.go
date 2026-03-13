package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type MigratorConfig struct {
	PostgresConfig postgresConfig
	LoggerConfig   loggerConfig  `yaml:"logger"`
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

func MustLoadConfig(configPath string) *MigratorConfig {
	cfg := new(MigratorConfig)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(err)
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic(err)
	}

	return cfg
}
