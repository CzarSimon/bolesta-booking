package main

import (
	"github.com/CzarSimon/bolesta-booking/backend/internal/api"
	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.GetConfig()
	api.Start(cfg)
}
