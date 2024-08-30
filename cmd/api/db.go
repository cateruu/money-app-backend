package main

import (
	"context"
	"log"
	"time"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getDBConfig(cfg config) *pgxpool.Config {
	var (
		maxOpenConns = cfg.db.maxOpenConns
		maxIdleTime  = cfg.db.maxIdleTime
	)

	dbConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = int32(maxOpenConns)
	dbConfig.MaxConnIdleTime = maxIdleTime

	dbConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxuuid.Register(c.TypeMap())
		return nil
	}

	return dbConfig
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	db, err := pgxpool.NewWithConfig(context.Background(), getDBConfig(cfg))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connection, err := db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	err = connection.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
