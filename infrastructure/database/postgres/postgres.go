package postgres

import (
	"construction_transport_server/config"
	"construction_transport_server/pkg/logger"
	"construction_transport_server/pkg/utils"
	"context"
	"log"
)

func Postgres(ctx context.Context) {
	cfg := config.LoadConfig().Db

	dbClient, err := New(ctx, cfg, &logger.SimpleLogger{}, &utils.NoopMetrics{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// 4. Ensure graceful shutdown
	defer dbClient.Close()
}
