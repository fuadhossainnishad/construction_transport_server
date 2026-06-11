package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"construction_transport_server/api/rest/v1"
	"construction_transport_server/api/rest/v1/delivery"
	"construction_transport_server/config"
	"construction_transport_server/infrastructure/cache/redis"
	"construction_transport_server/infrastructure/database/postgres"
	"construction_transport_server/infrastructure/messaging/rabbitmq"
	"construction_transport_server/internal/account/repository"
	accountUsecase "construction_transport_server/internal/account/usecase"
	"construction_transport_server/internal/auth/event"
	authRepository "construction_transport_server/internal/auth/repository"
	authUsecase "construction_transport_server/internal/auth/usecase"
	bookingRepository "construction_transport_server/internal/booking/repository"
	bookingUsecase "construction_transport_server/internal/booking/usecase"
	jobUsecase "construction_transport_server/internal/job/usecase"
	"construction_transport_server/internal/notification"
	vehicleRepository "construction_transport_server/internal/vehicle/repository"
	vehicleUsecase "construction_transport_server/internal/vehicle/usecase"
	"construction_transport_server/internal/websocket"
	"construction_transport_server/pkg/logger"
	"construction_transport_server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/client"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	cfg := config.LoadConfig()

	// PostgreSQL
	dbClient, err := postgres.New(ctx, cfg.Db, &logger.SimpleLogger{}, &utils.NoopMetrics{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer dbClient.Close()

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	otpStore := redis.NewOTPStore(redisClient)

	// RabbitMQ
	rabbitConn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("rabbitmq connection failed: %v", err)
	}
	defer rabbitConn.Close()
	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer rabbitCh.Close()
	// Declare queue
	_, err = rabbitCh.QueueDeclare("user.registered", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}
	rabbitPublisher := rabbitmq.NewPublisher(rabbitCh)

	// Event publisher
	eventPublisher := event.NewEventPublisher(rabbitPublisher)

	// Email consumer (runs in background)
	go notification.StartEmailConsumer(rabbitConn)

	// Stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	stripeClient := &client.API{}
	stripeClient.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	// JWT Manager
	jwtManager := utils.NewJWTManager(os.Getenv("JWT_SECRET"))

	// Repositories & Usecases
	authRepo := authRepository.NewAuthRepository(dbClient.Pool)
	hashFunc := &utils.BcryptPasswordHasher{}
	otpService := authUsecase.NewOTPService(otpStore, eventPublisher)
	refreshTokenRepo := authRepository.NewRefreshTokenRepository(dbClient.Pool)
	refreshUsecase := authUsecase.NewRefreshTokenUsecase(refreshTokenRepo)
	loginUsecase := authUsecase.NewLoginUseCase(authRepo, jwtManager, refreshUsecase, hashFunc)
	registerUsecase := authUsecase.NewRegisteredUsecase(authRepo, hashFunc, otpService, eventPublisher)
	authHandler := delivery.NewAuthHandler(registerUsecase, loginUsecase, otpService, authRepo, hashFunc, refreshUsecase)

	// Account (Profile)
	accountRepo := repository.NewAccountRepository(dbClient.Pool)
	accountUsecase := accountUsecase.NewAccountUsecase(accountRepo)
	accountHandler := delivery.NewAccountHandler(accountUsecase)

	// Vehicle
	vehicleRepo := vehicleRepository.NewVehicleRepository(dbClient.Pool)
	vehicleUsecase := vehicleUsecase.NewVehicleUsecase(vehicleRepo)
	vehicleHandler := delivery.NewVehicleHandler(vehicleUsecase)

	// Booking
	bookingRepo := bookingRepository.NewBookingRepository(dbClient.Pool)
	bookingUsecase := bookingUsecase.NewBookingUsecase(bookingRepo, vehicleRepo, eventPublisher)
	bookingHandler := delivery.NewBookingHandler(bookingUsecase)

	// WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Job
	jobUsecase := jobUsecase.NewJobUsecase(bookingRepo, wsHub, eventPublisher, stripeClient)
	jobHandler := delivery.NewJobHandler(jobUsecase)

	// Gin router
	router := gin.Default()
	v1.RegisterRoutes(router, authHandler, accountHandler, vehicleHandler, bookingHandler, jobHandler, jwtManager, wsHub)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 Server started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
	srv.Shutdown(context.Background())
}
