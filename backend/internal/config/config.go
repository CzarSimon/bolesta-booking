package config

import (
	"os"
	"strconv"

	"github.com/CzarSimon/httputil/dbutil"
	"github.com/CzarSimon/httputil/environ"
	"github.com/CzarSimon/httputil/logger"
)

var log = logger.GetDefaultLogger("internal/config").Sugar()

// Config application configuration.
type Config struct {
	DB                dbutil.Config
	MigrationsPath    string
	Port              string
	EnableCreateUsers bool
}

// GetConfig reads, parses and marshalls the applications configuration.
func GetConfig() Config {
	return Config{
		DB: dbutil.SqliteConfig{
			Name: environ.MustGet("DB_FILENAME"),
		},
		MigrationsPath:    environ.Get("MIGRATIONS_PATH", "/etc/bolesta-booking/backend/db/sqlite"),
		Port:              environ.Get("PORT", "8080"),
		EnableCreateUsers: getBoolEnvVar("ENABLE_CREATE_USERS", false),
	}
}

func getBoolEnvVar(key string, defaultValue bool) bool {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}

	val, err := strconv.ParseBool(str)
	if err != nil {
		log.Fatalf("failed to parse %s as boolean. Key=%s: %w", str, key, err)
	}

	return val
}
