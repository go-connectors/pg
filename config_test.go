package pg

import (
	"fmt"
	"testing"
)

var config = Config{
	Host:     "127.0.0.1",
	Port:     "5432",
	Database: "pcs",
	User:     "postgres",
	Password: "postgres",
	PoolSize: 10,
	Debug:    true,
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{"positive validation", config, false},
		{"empty host", Config{}, true},
		{"empty port", Config{Host: "localhost"}, true},
		{"empty database", Config{Host: "localhost", Port: "5432"}, true},
		{"empty user", Config{Host: "localhost", Port: "5432", Database: "pcs"}, true},
		{"empty password", Config{Host: "localhost", Port: "5432", Database: "pcs", User: "postgres"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validation error expected: %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestConfig_GetDSN(t *testing.T) {
	expected := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database,
	)

	if got := config.GetDSN(); got != expected {
		t.Errorf("dns expected: %v, got %v", expected, got)
	}
}
