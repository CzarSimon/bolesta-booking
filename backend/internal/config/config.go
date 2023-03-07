package config

import (
	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/httputil/dbutil"
	"github.com/CzarSimon/httputil/environ"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/CzarSimon/httputil/logger"
)

var log = logger.GetDefaultLogger("internal/config").Sugar()

// Config application configuration.
type Config struct {
	DB                dbutil.Config
	MigrationsPath    string
	Port              string
	EnableCreateUsers bool
	JWT               jwt.Credentials
	BookingRules      models.BookingRules
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
		BookingRules:      getBookingRules(),
	}
}

func getJWTCredentials() jwt.Credentials {
	return jwt.Credentials{
		Issuer: environ.Get("JWT_ISSUER", "bolesta-booking"),
		Secret: environ.MustGet("JWT_SECRET"),
	}
}

func getBookingRules() models.BookingRules {
	d := models.DefaultBookingRules()

	return models.BookingRules{
		MaxBookingLengthDays:  getIntEnvVar("BOOKING_RULE_MAX_LENGTH_IN_DAYS", d.MaxBookingLengthDays),
		MaxActiveBookings:     getIntEnvVar("BOOKING_RULE_MAX_ACTIVE_BOOKINGS", d.MaxActiveBookings),
		MustStartWithinMonths: getIntEnvVar("BOOKING_RULE_MUST_START_WITHIN_MONTHS", d.MustStartWithinMonths),
	}
}
