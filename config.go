package pg

import (
	"errors"
	"fmt"
)

// Config validation errors.
var (
	// ErrConfigValidation is general config validation error message.
	ErrConfigValidation = errors.New("pg config validation error")

	ErrEmptyHost     = fmt.Errorf("%w: host is empty", ErrConfigValidation)
	ErrEmptyPort     = fmt.Errorf("%w: port is empty", ErrConfigValidation)
	ErrEmptyDatabase = fmt.Errorf("%w: database is empty", ErrConfigValidation)
	ErrEmptyUser     = fmt.Errorf("%w: user is empty", ErrConfigValidation)
	ErrEmptyPassword = fmt.Errorf("%w: password is empty", ErrConfigValidation)
)

// Config contains config for PostgreSQL database.
type Config struct {
	Host         string            `yaml:"host"`
	Port         string            `yaml:"port"`
	Database     string            `yaml:"database"`
	User         string            `yaml:"username"`
	Password     string            `yaml:"password"`
	Notify       map[string]string `yaml:"notify"`
	Debug        bool              `yaml:"debug"`
	PoolSize     int               `yaml:"pool_size"`
	MaxIdleConns int               `yaml:"max_idle_conns"`
	MaxOpenConns int               `yaml:"max_open_conns"`
}

// Validate checks required fields and validates for allowed values.
func (cfg *Config) Validate() error {
	if cfg.Host == "" {
		return ErrEmptyHost
	}

	if cfg.Port == "" {
		return ErrEmptyPort
	}

	if cfg.Database == "" {
		return ErrEmptyDatabase
	}

	if cfg.User == "" {
		return ErrEmptyUser
	}

	if cfg.Password == "" {
		return ErrEmptyPassword
	}

	return nil
}

// GetDSN returns dsn connection string to postgres database.
func (cfg *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
