package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/CzarSimon/bolesta-booking/backend/internal/api/bookings"
	"github.com/CzarSimon/bolesta-booking/backend/internal/api/cabins"
	"github.com/CzarSimon/bolesta-booking/backend/internal/api/users"
	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/dbutil"
	"github.com/CzarSimon/httputil/logger"
	"github.com/gin-contrib/cors"
	"go.uber.org/zap"
)

var log = logger.GetDefaultLogger("internal/api")

func Start(cfg config.Config) {
	db := dbutil.MustConnect(cfg.DB)
	defer db.Close()

	err := dbutil.Upgrade(cfg.MigrationsPath, cfg.DB.Driver(), db)
	if err != nil {
		log.Panic("Failed to apply upgrade migratons", zap.Error(err))
	}

	r := httputil.NewRouter("bölesta-booking/backend", healthCheck(db))
	r.Use(cors.Default())

	userRepo := repository.NewUserRepository(db)
	cabinRepo := repository.NewCabinRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	userSvc := &service.UserService{
		UserRepo: userRepo,
	}
	cabinSvc := &service.CabinService{
		CabinRepo: cabinRepo,
	}
	bookingSvc := &service.BookingService{
		UserRepo:    userRepo,
		CabinRepo:   cabinRepo,
		BookingRepo: bookingRepo,
	}

	users.AttachController(userSvc, r, cfg)
	cabins.AttachController(cabinSvc, r)
	bookings.AttachController(bookingSvc, r)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Info("Starting bölesta-booking/backend", zap.String("port", cfg.Port))
	err = server.ListenAndServe()
	if err != nil {
		log.Error("Server stoped with an error", zap.Error(err))
	}
}

func healthCheck(db *sql.DB) httputil.HealthFunc {
	return func() error {
		return dbutil.Connected(db)
	}
}
