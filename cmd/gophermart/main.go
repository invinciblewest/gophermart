package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/invinciblewest/gophermart/internal/client/accrual"
	"github.com/invinciblewest/gophermart/internal/config"
	"github.com/invinciblewest/gophermart/internal/handler"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/repository/postgres"
	"github.com/invinciblewest/gophermart/internal/usecase/app"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	loadEnv()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err = logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.Fatal("failed to close database", zap.Error(err))
		}
	}(db)

	logger.Log.Info("attempting to connect to database...", zap.String("url", cfg.DatabaseURL))
	if err = db.Ping(); err != nil {
		logger.Log.Fatal("failed to ping database", zap.Error(err))
	}

	if err = runMigrations(db); err != nil {
		logger.Log.Fatal("failed to run migrations", zap.Error(err))
	}

	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress)

	repository := postgres.NewPGRepository(db)

	authUseCase := app.NewAuthUseCase(cfg.SecretKey)
	userUseCase := app.NewUserUseCase(repository, authUseCase)
	orderUseCase := app.NewOrderUseCase(repository)
	balanceUseCase := app.NewBalanceUseCase(repository, repository)

	accrualProcessor := app.NewAccrualProcessor(repository, accrualClient)

	go accrualProcessor.Run(ctx, cfg.UpdateInterval, cfg.WorkerCount)

	router := handler.NewRouter(
		handler.NewHandler(userUseCase, orderUseCase, balanceUseCase),
		authUseCase,
	)

	if err = runHTTPServer(ctx, cfg.RunAddress, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("HTTP server error", zap.Error(err))
	}
}

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return nil
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func runHTTPServer(ctx context.Context, address string, handler http.Handler) error {
	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Log.Info("server is shutting down...")
		if err := server.Shutdown(ctxWithTimeout); err != nil {
			logger.Log.Fatal("server shutdown error", zap.Error(err))
		}
	}()

	logger.Log.Info("server is starting", zap.String("address", address))

	return server.ListenAndServe()
}

func loadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err = godotenv.Load(); err != nil {
			log.Println("error loading .env file:", err)
		}
	}
}
