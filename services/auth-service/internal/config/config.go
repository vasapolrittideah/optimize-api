package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

// AuthServiceConfig contains the configuration for the auth service.
type AuthServiceConfig struct {
	Environment string `env:"ENVIRONMENT"`
	Name        string `env:"SERVICE_NAME"`
	Address     string `env:"SERVICE_ADDRESS"`
	Token       TokenConfig
}

// TokenConfig contains the configuration for JWT tokens.
type TokenConfig struct {
	AccessTokenSecret     string        `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret    string        `env:"REFRESH_TOKEN_SECRET"`
	AccessTokenExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRES_IN"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRES_IN"`
	Issuer                string        `env:"TOKEN_ISSUER"`
}

// NewAuthServiceConfig creates a new AuthServiceConfig instance from environment variables.
func NewAuthServiceConfig(logger *zerolog.Logger) *AuthServiceConfig {
	cfg, err := env.ParseAs[AuthServiceConfig]()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse environment variables")
	}

	return &cfg
}
