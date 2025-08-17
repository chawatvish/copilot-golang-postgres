package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
		wantErr  bool
	}{
		{
			name: "success_with_defaults",
			envVars: map[string]string{
				"SERVER_PORT": "8080",
			},
			expected: &Config{
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "password",
					Name:     "gin_app",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port:    "8080",
					GinMode: "debug",
				},
				JWT: JWTConfig{
					Secret:     "your-secret-key-change-this-in-production",
					ExpireHour: 24,
				},
			},
			wantErr: false,
		},
		{
			name: "success_with_custom_values",
			envVars: map[string]string{
				"DB_HOST":        "custom-host",
				"DB_PORT":        "3306",
				"DB_USER":        "custom-user",
				"DB_PASSWORD":    "custom-pass",
				"DB_NAME":        "custom-db",
				"DB_SSLMODE":     "require",
				"SERVER_PORT":    "3000",
				"GIN_MODE":       "release",
				"JWT_SECRET":     "custom-secret",
				"JWT_EXPIRE_HOUR": "12",
			},
			expected: &Config{
				Database: DatabaseConfig{
					Host:     "custom-host",
					Port:     "3306",
					User:     "custom-user",
					Password: "custom-pass",
					Name:     "custom-db",
					SSLMode:  "require",
				},
				Server: ServerConfig{
					Port:    "3000",
					GinMode: "release",
				},
				JWT: JWTConfig{
					Secret:     "custom-secret",
					ExpireHour: 12,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_jwt_expire_hour",
			envVars: map[string]string{
				"JWT_EXPIRE_HOUR": "invalid-number",
			},
			expected: &Config{
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "password",
					Name:     "gin_app",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port:    "8080",
					GinMode: "debug",
				},
				JWT: JWTConfig{
					Secret:     "your-secret-key-change-this-in-production",
					ExpireHour: 24, // Should fallback to default
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original environment variables
			originalEnvVars := make(map[string]string)
			for key := range tt.envVars {
				originalEnvVars[key] = os.Getenv(key)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Restore original environment variables after test
			defer func() {
				for key, originalValue := range originalEnvVars {
					if originalValue == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, originalValue)
					}
				}
			}()

			config, err := Load()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, config)
			}
		})
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	dbConfig := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
	}
	
	expectedDSN := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	dsn := dbConfig.GetDSN()
	assert.Equal(t, expectedDSN, dsn)
}

func TestJWTConfig_GetJWTExpiry(t *testing.T) {
	jwtConfig := JWTConfig{
		ExpireHour: 12,
	}
	
	expected := 12 * time.Hour
	expiry := jwtConfig.GetJWTExpiry()
	assert.Equal(t, expected, expiry)
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "env_var_exists",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "test-value",
			expected:     "test-value",
		},
		{
			name:         "env_var_not_exists",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "env_var_empty_string",
			key:          "EMPTY_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original value
			originalValue := os.Getenv(tt.key)
			defer func() {
				if originalValue == "" {
					os.Unsetenv(tt.key)
				} else {
					os.Setenv(tt.key, originalValue)
				}
			}()

			// Set test environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
