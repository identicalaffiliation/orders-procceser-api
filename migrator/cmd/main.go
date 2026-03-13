package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/identicalaffiliation/orders-procceser-api/migrator/internal/config"
	"github.com/identicalaffiliation/orders-procceser-api/migrator/internal/logger"
	"github.com/identicalaffiliation/orders-procceser-api/migrator/internal/storage"
)

const (
	CREATE_ORDERS_TABLE string = `
		CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY NOT NULL,
			status varchar(20) NOT NULL,
			total_price NUMERIC(10, 2) check (total_price > 0) NOT NULL,
			total_quantity INT check (total_quantity > 0) NOT NULL,
			created TIMESTAMPTZ DEFAULT NOW(),
			updated TIMESTAMPTZ DEFAULT NOW()
		);
	`

	CREATE_ITEMS_TABLE string = `
		CREATE TABLE IF NOT EXISTS items (
			id UUID PRIMARY KEY NOT NULL,
			order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			price NUMERIC(10, 2) check (price > 0) NOT NULL,
			quantity INT check (quantity > 0) NOT NULL,
			created TIMESTAMPTZ DEFAULT NOW(),
			updated TIMESTAMPTZ DEFAULT NOW()
		);
	`
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg := config.MustLoadConfig(configPath)
	slogger := logger.NewLogger(cfg)

	psql, err := storage.NewConnect(cfg)
	if err != nil {
		slogger.Error("connect to postgres", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := psql.Close(); err != nil {
			slogger.Error("close postgres", "error", err)
			os.Exit(1)
		}
	}()

	ctx, cancel := createContext(cfg)
	defer cancel()

	if err := initMigrations(ctx, psql, slogger); err != nil {
		slogger.Error("init migrations", "error", err)
		os.Exit(1)
	}

	slogger.Debug("migrations is added")
	os.Exit(0)
}

func initMigrations(ctx context.Context, psql *storage.Postgres, logger logger.Logger) error {
	tx, err := psql.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			logger.Error("rollback tx", "error", err)
			os.Exit(1)
		}
	}()

	if _, err := tx.ExecContext(ctx, CREATE_ORDERS_TABLE); err != nil {
		logger.Error("exec orders", "error", err)
		os.Exit(1)
	}

	if _, err := tx.ExecContext(ctx, CREATE_ITEMS_TABLE); err != nil {
		logger.Error("exec items", "error", err)
		os.Exit(1)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit tx", "error", err)
		os.Exit(1)
	}

	return nil
}

func createContext(cfg *config.MigratorConfig) (context.Context, context.CancelFunc) {
	pCtx := context.Background()
	return context.WithTimeout(pCtx, cfg.Timeout)
}
