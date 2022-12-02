package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v2"
)

// SecurityConfig contains the configuration of security parameters.
type SecurityConfig struct {
	HashIterations  uint32 `yaml:"hash_iterations" env:"HASH_ITERATIONS" env-default:"3"`
	HashMemory      uint32 `yaml:"hash_memory" env:"HASH_MEMORY" env-default:"64"` // KiB
	HashParallelism uint8  `yaml:"hash_parallelism" env:"HASH_PARALLELISM" env-default:"2"`
	SaltLength      uint32 `yaml:"salt_length" env:"SALT_LENGTH" env-default:"16"`
	KeyLength       uint32 `yaml:"key_length" env:"KEY_LENGTH" env-default:"32"`
	// AtExpiresIn is responsible for the expiration time of the `jwt access token` in minutes.
	AtExpiresIn int `yaml:"at_expires_in" env:"AT_EXPIRES_IN" env-default:"15"`
	// RtExpiresIn is responsible for the expiration time of the `jwt refresh token` in minutes.
	RtExpiresIn   int    `yaml:"rt_expires_in" env:"RT_EXPIRES_IN" env-default:"10080"`
	AccessSecret  string `yaml:"access_secret" env:"ACCESS_SECRET" env-default:"tNJE2FjDGwWWBHZCx2Kz8NH95r7pxmwX"`
	RefreshSecret string `yaml:"refresh_secret" env:"REFRESH_SECRET" env-default:"JdAjGzK84R7okdKumM3MBEt28jnJC2cG"`
}

// WebConfig contains the web server configuration.
type WebConfig struct {
	Network       string `yaml:"network" env:"NETWORK" env-default:"tcp"`
	ServerAddress string `yaml:"server_address" env:"SERVER_ADDRESS" env-default:":8080"`
	TLSCertPath   string `yaml:"tls_cert_path" env:"TLS_CERT_PATH" env-default:"./config/server-cert.pem"`
	TLSKeyPath    string `yaml:"tls_key_path" env:"TLS_KEY_PATH" env-default:"./config/server-key.pem"`
	BasePath      string `yaml:"base_path" env-default:"/" env:"BASE_PATH"`
	EnableHTTPS   bool   `yaml:"enable_https" env-default:"true" env:"ENABLE_HTTPS"`
}

// DBConfig contains the database configuration.
type DBConfig struct {
	Type           string `yaml:"type" env-default:"sqlite3" env:"DB_TYPE"`
	DSN            string `yaml:"dsn" env-default:"db.sqlite3" env:"DB_DSN"`
	MigrationsPath string `yaml:"migrations_path,omitempty" env:"MIGRATIONS_PATH"`
}

// Version contains the Version information.
type Version struct {
	Version   string
	Commit    string
	BuildDate string
}

// ServerCfg contains the configuration of the application.
type ServerCfg struct {
	Web      WebConfig
	DB       DBConfig
	Version  Version
	Security SecurityConfig
}

// String returns version, build date and commit id.
func (v Version) String() string {
	return fmt.Sprintf("\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n",
		v.Version,
		v.BuildDate,
		v.Commit,
	)
}

const (
	srvConfigPath                = "./config/config.yml"
	ClientConfigPath             = "./config.yml"
	NA                           = "N/A"
	FilePermission   os.FileMode = 0o600
	MaxFileSize                  = 100 * 1024 * 1024
)

// NewServerConfig creates a new ServerCfg.
func NewServerConfig() (*ServerCfg, error) {
	cfg := ServerCfg{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	if _, err = os.ReadFile(srvConfigPath); err == nil {
		err = cleanenv.ReadConfig(srvConfigPath, &cfg)
		if err != nil {
			return nil, err
		}
	}

	cfg.Security.HashMemory *= 1024

	cfg.Web.TLSCertPath, _ = filepath.Abs(cfg.Web.TLSCertPath)
	cfg.Web.TLSKeyPath, _ = filepath.Abs(cfg.Web.TLSKeyPath)

	buildVersion(&cfg.Version)

	return &cfg, nil
}

// ClientCfg contains the configuration of the cli application.
type ClientCfg struct {
	DB         DBConfig       `yaml:"db,omitempty"`
	Version    Version        `yaml:"-"`
	Security   SecurityConfig `yaml:"security,omitempty"`
	SkipVerify bool           `yaml:"skip_verify" ENV:"skip_verify" env-default:"1"`
	API        API            `yaml:"api_client"`
}

// SaveConfig saves the current configuration to a file.
func (c *ClientCfg) SaveConfig() error {
	configToStore := &c

	(*configToStore).Security = SecurityConfig{}

	data, err := yaml.Marshal(&configToStore)
	if err != nil {
		return err
	}

	return os.WriteFile(ClientConfigPath, data, FilePermission)
}

// API contains the configuration for communicating with the
// remote server.
type API struct {
	Remote string `yaml:"remote" ENV:"REMOTE"`
	AT     string `yaml:"at" ENV:"ACCESS_TOKEN"`
	RT     string `yaml:"rt" ENV:"REFRESH_TOKEN"`
}

// NewClientConfig creates a new ClientConfig.
func NewClientConfig() (*ClientCfg, error) {
	cfg := ClientCfg{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	if _, err = os.ReadFile(ClientConfigPath); err == nil {
		err = cleanenv.ReadConfig(ClientConfigPath, &cfg)
		if err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				return nil, err
			}
		}
	}

	buildVersion(&cfg.Version)

	return &cfg, nil
}

func buildVersion(cfg *Version) {
	if cfg.Version == "" {
		cfg.Version = "N/A"
	}
	if cfg.BuildDate == "" {
		cfg.BuildDate = "N/A"
	}

	cfg.Commit = commit()
}
