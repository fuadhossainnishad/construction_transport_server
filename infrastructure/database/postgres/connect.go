package postgres

import (
	"context"
	"fmt"
	"time"

	"construction_transport_server/config"
	"construction_transport_server/pkg/logger"
	"construction_transport_server/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectWithRetry(
	ctx context.Context,
	cfg *pgxpool.Config,
	dbCfg config.DBConfig,
	logger logger.Logger,
	metrics utils.Metrics,
) (*pgxpool.Pool, error) {

	retry := NewExponentialBackoff(5, time.Second, 5*time.Second)

	var pool *pgxpool.Pool
	var err error

	for retry.HasNext() {

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		pool, err = pgxpool.NewWithConfig(attemptCtx, cfg)
		if err == nil {
			pingErr := pool.Ping(attemptCtx)
			if pingErr == nil {

				cancel()

				logger.Info("[DB] connected",
					"host", dbCfg.Host,
					"db", dbCfg.DBName,
				)

				metrics.IncDBSuccess()
				return pool, nil
			}

			pool.Close()
			err = pingErr
		}

		cancel()

		delay := retry.Next()

		logger.Error("[DB] retrying",
			"attempt", retry.attempts,
			"error", err,
			"delay", delay,
		)

		metrics.IncDBRetry()

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	metrics.IncDBFailure()
	return nil, fmt.Errorf("postgres connection failed after retries: %w", err)
}
