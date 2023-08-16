package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"project/api"
	"time"
)

//go:generate mockgen -source=.\irepository.go -destination=mocks\repository_mock.go

type IRepository interface {
	GetItem(ctx context.Context, uuid uuid.UUID) (api.Item, error)
	GetItems(ctx context.Context) ([]ItemRow, error)
	AddItem(ctx context.Context, uuid uuid.UUID, item api.Item)
}

type Config interface {
	GetDSN() string
	GetMaxOpenConn() int
	GetMaxIdleConn() int
	GetConnMaxLifetime() time.Duration
	GetConnMaxIdleTime() time.Duration
}

func ConnectPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.GetDSN())

	db.SetMaxOpenConns(cfg.GetMaxOpenConn())
	db.SetMaxIdleConns(cfg.GetMaxIdleConn())
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())
	db.SetConnMaxIdleTime(cfg.GetConnMaxIdleTime())

	if err != nil {
		return nil, err
	}
	return db, nil
}
