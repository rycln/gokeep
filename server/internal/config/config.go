// Package config provides centralized application configuration management.
package config

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/caarlos0/env/v11"
	"github.com/rycln/gokeep/server/internal/logger"
)

// Config default values
const (
	defaultGRPCPort  = ":50051"
	defaultTimeout   = time.Duration(2) * time.Minute
	defaultKeyLength = 32
	defaultLogLevel  = "debug"
)

var errEmptyCfgFilepath = errors.New("empty cfg file path")

// Cfg contains all application configuration parameters.
//
// The structure supports loading from multiple sources:
// - Environment variables (primary)
// - Command-line flags (secondary)
// - Config file (tertiary)
// - Default values (fallback)
//
// Tags specify the corresponding environment variable names.
type Cfg struct {
	// DatabaseDsn specifies database connection string
	DatabaseDsn string `json:"database_dsn" env:"DATABASE_DSN"`

	// Key contains JWT signing key (min 32 bytes recommended)
	Key string `json:"jwt_key" env:"JWT_KEY"`

	// LogLevel sets logging verbosity (debug|info|warn|error)
	LogLevel string `json:"log_level" env:"LOG_LEVEL"`

	// GRPCPort defines port for gRPC endpoints
	GRPCPort string `json:"grpc_port" env:"GRPC_PORT"`

	// CfgFileName specifies configuration file name
	CfgFileName string `json:"-" env:"CONFIG"`

	// CertFileName specifies cert file name
	CertFileName string `json:"cert" env:"CERT"`

	// CertFileName specifies cert key file name
	CertKeyFileName string `json:"cert_key" env:"CERT_KEY"`

	// Timeout defines default network operation timeout
	Timeout time.Duration `json:"timeout_dur" env:"TIMEOUT_DUR"`
}

// ConfigBuilder implements builder pattern for Cfg.
type ConfigBuilder struct {
	cfg *Cfg
	err error
}

// NewConfigBuilder creates a new configuration builder with default values.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		cfg: &Cfg{
			Timeout:  defaultTimeout,
			LogLevel: defaultLogLevel,
			GRPCPort: defaultGRPCPort,
		},
		err: nil,
	}
}

// WithConfigFile load configuration values from specified file.
func (b *ConfigBuilder) WithConfigFile() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	if b.cfg.CfgFileName != "" {
		err := getCfgFromFile(b.cfg.CfgFileName, b.cfg)
		if err != nil {
			b.cfg = nil
			b.err = fmt.Errorf("can't open cfg file: %v", err)
			return b
		}
	} else {
		b.cfg = nil
		b.err = errEmptyCfgFilepath
		return b
	}

	return b
}

func getCfgFromFile(fname string, cfg *Cfg) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return err
	}

	return nil
}

// WithFlagParsing parses command-line flags into configuration.
func (b *ConfigBuilder) WithFlagParsing() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	flag.StringVarP(&b.cfg.DatabaseDsn, "d", "d", b.cfg.DatabaseDsn, "Database connection address")
	flag.DurationVarP(&b.cfg.Timeout, "t", "t", b.cfg.Timeout, "Timeout duration in seconds")
	flag.StringVarP(&b.cfg.Key, "k", "k", b.cfg.Key, "Key for jwt autorization")
	flag.StringVarP(&b.cfg.LogLevel, "l", "l", b.cfg.LogLevel, "Logger level")
	flag.StringVarP(&b.cfg.GRPCPort, "g", "g", b.cfg.GRPCPort, "gRPC port")
	flag.StringVarP(&b.cfg.CfgFileName, "config", "c", b.cfg.CfgFileName, "Path to config file")
	flag.StringVar(&b.cfg.CertFileName, "tls-cert", b.cfg.CertFileName, "Path to cert file")
	flag.StringVar(&b.cfg.CertKeyFileName, "tls-key", b.cfg.CertKeyFileName, "Path to cert key file")
	flag.Parse()

	return b
}

// WithEnvParsing loads environment variables into configuration.
func (b *ConfigBuilder) WithEnvParsing() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	err := env.Parse(b.cfg)
	if err != nil {
		b.cfg = nil
		b.err = fmt.Errorf("can't parse env vars: %v", err)
		return b
	}

	return b
}

// WithDefaultJWTKey sets default jwt key.
func (b *ConfigBuilder) WithDefaultJWTKey() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	if b.cfg.Key == "" {
		key, err := generateKey(defaultKeyLength)
		if err != nil {
			b.cfg = nil
			b.err = fmt.Errorf("can't generate jwt key: %v", err)
			return b
		}
		b.cfg.Key = key
		logger.Log.Warn("Default JWT key used!")
	}

	return b
}

func generateKey(n int) (string, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return string(key), nil
}

// Build finalizes configuration.
func (b *ConfigBuilder) Build() (*Cfg, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.cfg, nil
}
