package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	modeDevelopment = "development"
	modeProduction  = "production"
)

var (
	errEndpointRequired        = errors.New("Coda API endpoint is required")
	errEndpointInvalid         = errors.New("Coda API endpoint is invalid")
	errDatabaseRequired        = errors.New("Database credentials are required")
	errSyncIntervalRequired    = errors.New("Sync interval is required")
	errSyncIntervalInvalid     = errors.New("Sync interval is invalid")
	errCleanupIntervalRequired = errors.New("Cleanup interval is required")
	errCleanupIntervalInvalid  = errors.New("Cleanup interval is invalid")
)

// Config holds the configration data
type Config struct {
	AppEnv           string `json:"app_env" envconfig:"APP_ENV" default:"development"`
	MinaEndpoint     string `json:"mina_endpoint" envconfig:"MINA_ENDPOINT"`
	ArchiveEndpoint  string `json:"archive_endpoint" envconfig:"ARCHIVE_ENDPOINT"`
	StaketabEndpoint string `json:"staketab_endpoint" envconfig:"STAKETAB_ENDPOINT" default:"https://api.staketab.com/mina/get_all_providers"`
	GenesisFile      string `json:"genesis_file" envconfig:"GENESIS_FILE"`
	IdentityFile     string `json:"identity_file" envconfig:"IDENTITY_FILE"`
	ServerAddr       string `json:"server_addr" envconfig:"SERVER_ADDR" default:"0.0.0.0"`
	ServerPort       int    `json:"server_port" envconfig:"SERVER_PORT" default:"8080"`
	SyncInterval     string `json:"sync_interval" envconfig:"SYNC_INTERVAL" default:"60s"`
	CleanupInterval  string `json:"cleanup_interval" envconfig:"CLEANUP_INTERVAL" default:"10m"`
	CleanupThreshold int    `json:"cleanup_threshold" envconfig:"CLEANUP_THRESHOLD" default:"1000"`
	DatabaseURL      string `json:"database_url" envconfig:"DATABASE_URL"`
	DumpDir          string `json:"dump_dir" envconfig:"DUMP_DIR"`
	LogLevel         string `json:"log_level" envconfig:"LOG_LEVEL" default:"info"`
	LogFormat        string `json:"log_format" envconfig:"LOG_FORMAT" default:"text"`
	RollbarToken     string `json:"rollbar_token" envconfig:"ROLLBAR_TOKEN"`
	RollbarNamespace string `json:"rollbar_namespace" envconfig:"ROLLBAR_NAMESPACE"`

	HistoricalLimit uint `json:"historical_limit" envconfig:"HISTORICAL_LIMIT" default:"290"`

	syncDuration    time.Duration
	cleanupDuration time.Duration
}

// Validate returns an error if config is invalid
func (c *Config) Validate() error {
	if c.MinaEndpoint == "" {
		return errEndpointRequired
	}
	codaURL, err := url.Parse(c.MinaEndpoint)
	if err != nil {
		return errEndpointInvalid
	}
	if !strings.Contains(codaURL.Path, "graphql") {
		return errEndpointInvalid
	}

	if c.DatabaseURL == "" {
		return errDatabaseRequired
	}

	if c.SyncInterval == "" {
		return errSyncIntervalRequired
	}
	d, err := time.ParseDuration(c.SyncInterval)
	if err != nil {
		return errSyncIntervalInvalid
	}
	c.syncDuration = d

	if c.CleanupInterval == "" {
		return errCleanupIntervalRequired
	}
	d, err = time.ParseDuration(c.CleanupInterval)
	if err != nil {
		return errCleanupIntervalInvalid
	}
	c.cleanupDuration = d

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

// SyncDuration returns the parsed duration for the sync pipeline
func (c *Config) SyncDuration() time.Duration {
	return c.syncDuration
}

// CleanupDuration returns the parsed duration for the cleanup pipeline
func (c *Config) CleanupDuration() time.Duration {
	return c.cleanupDuration
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
