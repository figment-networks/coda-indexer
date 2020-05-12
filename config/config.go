package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
)

const (
	modeDevelopment = "development"
	modeProduction  = "production"
)

var (
	errEndpointRequired = errors.New("Coda API endpoint is required")
	errDatabaseRequired = errors.New("Database credentials are required")
)

// Config holds the configration data
type Config struct {
	AppEnv           string `json:"app_env" envconfig:"APP_ENV" default:"development"`
	CodaEndpoint     string `json:"coda_endpoint" envconfig:"CODA_ENDPOINT"`
	ServerAddr       string `json:"server_addr" envconfig:"SERVER_ADDR" default:"0.0.0.0"`
	ServerPort       int    `json:"server_port" envconfig:"SERVER_PORT" default:"8081"`
	FirstBlockHeight int    `json:"first_block_height" envconfig:"FIRST_BLOCK_HEIGHT" default:"1"`
	SyncInterval     string `json:"sync_interval" envconfig:"SYNC_INTERVAL" default:"10s"`
	CleanupInterval  string `json:"cleanup_interval" envconfig:"CLEANUP_INTERVAL" default:"10min"`
	CleanupThreshold int    `json:"cleanup_threshold" envconfig:"CLEANUP_THRESHOLD"`
	DatabaseURL      string `json:"database_url" envconfig:"DATABASE_URL"`
	Debug            bool   `json:"debug" envconfig:"DEBUG"`
}

// Validate returns an error if config is invalid
func (c *Config) Validate() error {
	if c.CodaEndpoint == "" {
		return errEndpointRequired
	}
	if c.DatabaseURL == "" {
		return errDatabaseRequired
	}
	return nil
}

// IsDevelopment returns true if app is in dev mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == modeDevelopment
}

// IsProduction returns true if app is in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == modeProduction
}

// ListenAddr returns a full listen address and port
func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
}

// New returns a new config
func New() *Config {
	return &Config{}
}

// FromFile reads the config from a file
func FromFile(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// FromEnv reads the config from environment variables
func FromEnv(config *Config) error {
	return envconfig.Process("", config)
}
