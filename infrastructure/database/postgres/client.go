package postgres

import (
	"construction_transport_server/config"
	"construction_transport_server/pkg/logger"
	"construction_transport_server/pkg/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	Pool    *pgxpool.Pool
	logger  logger.Logger
	metrics utils.Metrics
}

func New(
	ctx context.Context,
	cfg config.DBConfig,
	logger logger.Logger,
	metrics utils.Metrics,
) (*Client, error) {

	poolConfig, err := BuildPoolConfig(cfg)
	if err != nil {
		return nil, err
	}

	pool, err := ConnectWithRetry(ctx, poolConfig, cfg, logger, metrics)
	if err != nil {
		return nil, err
	}

	return &Client{
		Pool:    pool,
		logger:  logger,
		metrics: metrics,
	}, nil
}
