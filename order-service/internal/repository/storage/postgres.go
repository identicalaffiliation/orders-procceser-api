package storage

import (
	"fmt"

	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	POSTGRES_DRIVER_NAME string = "postgres"
)

type Postgres struct {
	DB *sqlx.DB
}

func NewConnect(cfg *config.MigratorConfig) (*Postgres, error) {
	psql, err := sqlx.Open(POSTGRES_DRIVER_NAME, getDataSourceName(cfg))
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	if err := psql.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Postgres{DB: psql}, nil
}

func getDataSourceName(cfg *config.MigratorConfig) string {
	return fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresConfig.Port, cfg.PostgresConfig.Host, cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password, cfg.PostgresConfig.DBname,
		cfg.PostgresConfig.SSLmode)
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}
