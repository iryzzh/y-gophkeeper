package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

// SecurityConfig contains the configuration of security parameters.
type SecurityConfig struct {
	HashIterations  uint32 `yaml:"hash_iterations" env:"HASH_ITERATIONS" default:"3"`
	HashMemory      uint32 `yaml:"hash_memory" env:"HASH_MEMORY" default:"64"` // KiB
	HashParallelism uint8  `yaml:"hash_parallelism" env:"HASH_PARALLELISM" default:"2"`
	SaltLength      uint32 `yaml:"salt_length" env:"SALT_LENGTH" default:"16"`
	KeyLength       uint32 `yaml:"key_length" env:"KEY_LENGTH" default:"32"`
	// AtExpiresIn is responsible for the expiration time of the `jwt access token` in minutes.
	AtExpiresIn int `yaml:"at_expires_in" env:"AT_EXPIRES_IN" default:"15"`
	// RtExpiresIn is responsible for the expiration time of the `jwt refresh token` in minutes.
	RtExpiresIn   int    `yaml:"rt_expires_in" env:"RT_EXPIRES_IN" default:"10080"`
	AccessSecret  string `yaml:"access_secret" env:"ACCESS_SECRET" default:"tNJE2FjDGwWWBHZCx2Kz8NH95r7pxmwX"`
	RefreshSecret string `yaml:"refresh_secret" env:"REFRESH_SECRET" default:"JdAjGzK84R7okdKumM3MBEt28jnJC2cG"`
}

// WebServerConfig contains the web server configuration.
type WebServerConfig struct {
	Network       string `yaml:"network" env:"NETWORK" default:"tcp"`
	ServerAddress string `yaml:"server_address" env:"SERVER_ADDRESS" default:":8080"`
	TLSCertPath   string `yaml:"tls_cert_path" env:"TLS_CERT_PATH" default:"./config/server-cert.pem"`
	TLSKeyPath    string `yaml:"tls_key_path" env:"TLS_KEY_PATH" default:"./config/server-key.pem"`
	BasePath      string `yaml:"base_path" default:"/" env:"BASE_PATH"`
	EnableHTTPS   bool   `yaml:"enable_https" default:"true" env:"ENABLE_HTTPS"`
}

// DBConfig contains the database configuration.
type DBConfig struct {
	Type           string `yaml:"db_type" default:"sqlite3" env:"DB_TYPE"`
	DSN            string `yaml:"db_dsn" default:"db.sqlite3" env:"DB_DSN"`
	MigrationsPath string `yaml:"migrations_path" env:"MIGRATIONS_PATH"`
}

// Version contains the Version information.
type Version struct {
	Version   string
	Commit    string
	BuildDate string
}

// Config contains the configuration of the application.
type Config struct {
	WebServer WebServerConfig
	DB        DBConfig
	Version   Version
	Security  SecurityConfig
}

// GetVersion returns version, build date and commit id.
func (c *Config) GetVersion() string {
	if c.Version.Version == "" {
		c.Version.Version = "N/A"
	}
	if c.Version.BuildDate == "" {
		c.Version.BuildDate = "N/A"
	}

	c.Version.Commit = commit()

	return fmt.Sprintf("\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n",
		c.Version.Version,
		c.Version.BuildDate,
		c.Version.Commit,
	)
}

// NewConfig creates a new Config.
func NewConfig() (*Config, error) {
	config := &Config{}

	if err := configor.New(&configor.Config{}).Load(config); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	config.Security.HashMemory *= 1024

	config.WebServer.TLSCertPath, _ = filepath.Abs(config.WebServer.TLSCertPath)
	config.WebServer.TLSKeyPath, _ = filepath.Abs(config.WebServer.TLSKeyPath)

	return config, nil
}

// mustReadConfigLocation reads the config file from the given path.
func mustReadConfigLocation(cFlag *string) string {
	var err error
	*cFlag, err = filepath.Abs(*cFlag)
	if err != nil {
		panic(fmt.Sprintf("failed to get an absolute representation of path: %s: %v", *cFlag, err))
	}
	if _, err := os.Stat(*cFlag); err != nil {
		panic(fmt.Sprintf("config file not found: %s", *cFlag))
	}

	return *cFlag
}
