package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config ...
type Config struct {
	Addr                         string        `envconfig:"ADDR"`
	ServerSecret                 string        `envconfig:"SERVER_SECRET"`
	ServerListenersLimit         int64         `envconfig:"SERVER_LISTENERS_LIMIT"`
	Timeout                      time.Duration `envconfig:"TIMEOUT"`
	HashcashChallengeExpDuration time.Duration `envconfig:"HASHCASH_CHALLENGE_EXP_DURATION"`
	HashcashBits                 int32         `envconfig:"HASHCASH_BITS"`
	ClientMaxIterations          int64         `envconfig:"CLIENT_MAX_ITERATIONS"`
}

// Load load config
func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("WARN can't load .env error: %s", err)
	}

	config := new(Config)

	err = envconfig.Process("", config)
	if err != nil {
		return nil, fmt.Errorf("config processing error: %w", err)
	}

	return config, nil
}
