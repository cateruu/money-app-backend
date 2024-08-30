package main

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/pkg/logger"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type app struct {
	config config
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 1337, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|prod)")

	flag.StringVar(&cfg.db.dsn, "dsn", "", "Database connection DSN string")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	models := data.NewModels(db)

	app := app{
		config: cfg,
	}

	err = app.serve(&models)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}
}
