package main

import (
	"construction_transport_server/config"
	"construction_transport_server/infrastructure/database/postgres"
	"construction_transport_server/pkg/logger"
	"construction_transport_server/pkg/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()
	dbClient, err := postgres.New(ctx, cfg.Db, &logger.SimpleLogger{}, &utils.NoopMetrics{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// ✅ close ONLY when app exits
	defer dbClient.Close()

	// postgres.Postgres(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server Started Successfully on port", port)
	log.Println("🚀 Server Started Successfully on port", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello brother!")
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
