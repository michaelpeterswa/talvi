package cockroach

import (
	"context"
	"errors"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPgxFailedToParseConfig = errors.New("pgx failed to parse config")
	ErrPgxFailedToConnect     = errors.New("pgx failed to connect")
)

type CockroachClient struct {
	Client *pgxpool.Pool
}

func NewCockroachClient(connectionString string) (*CockroachClient, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, ErrPgxFailedToParseConfig
	}

	config.ConnConfig.RuntimeParams["application_name"] = "$ talvi_backend"
	config.ConnConfig.Tracer = otelpgx.NewTracer()

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &CockroachClient{Client: conn}, nil
}
