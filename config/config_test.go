package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Ensure env vars are not set
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("JWT_SECRET")

	cfg := LoadConfig()

	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "postgres", cfg.Database.Password)
	assert.Equal(t, "footballpaul.db", cfg.Database.DBName)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "your-secret-key-change-in-production", cfg.JWT.Secret)
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "admin")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("JWT_SECRET", "my-super-secret")

	cfg := LoadConfig()

	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "admin", cfg.Database.User)
	assert.Equal(t, "3000", cfg.Server.Port)
	assert.Equal(t, "my-super-secret", cfg.JWT.Secret)

	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("JWT_SECRET")
}

func TestDatabaseConfigGetDSN(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "secret",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	dsn := cfg.GetDSN()
	assert.Contains(t, dsn, "localhost")
	assert.Contains(t, dsn, "5432")
	assert.Contains(t, dsn, "postgres")
	assert.Contains(t, dsn, "secret")
	assert.Contains(t, dsn, "testdb")
	assert.Contains(t, dsn, "disable")
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "custom-value")
	assert.Equal(t, "custom-value", getEnv("TEST_VAR", "default"))
	os.Unsetenv("TEST_VAR")
	assert.Equal(t, "default", getEnv("TEST_VAR", "default"))
}
