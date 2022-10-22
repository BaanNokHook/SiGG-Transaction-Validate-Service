package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		Log
		RMQ
		PusherBeam
	}

	// App -.
	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	//RMQ -.
	RMQ struct {
		ClientExchange string `env-required:"true" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"true" env:"RMQ_URL"`
	}

	PusherBeam struct {
		InstanceId string `env-required:"true" env:"PUSHER_BEAM_INSTANCE_ID"`
		SecretKey  string `env-required:"true" env:"PUSHER_BEAM_SECRET_KEY"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(".env"); err == nil {
		err = cleanenv.ReadConfig(".env", cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
