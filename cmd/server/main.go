// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"canvas/messaging"
	"canvas/server"
	"canvas/storage"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maragudk/env"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// release is set through the linker at build time, generally from a git sha.
// passed during build process
// Used for logging and error reporting.
var release string

func main() {
	os.Exit(start())
}

func start() int {
	_ = env.Load()

	log, err := createLogger(env.GetStringOrDefault("LOG_ENV", "development"))
	if err != nil {
		fmt.Println("Error setting up logger:", err)
		return 1
	}

	log = log.With(zap.String("release", release))

	defer func() {
		// If we cannot sync, there's probably something wrong with outputting logs,
		// so we probably cannot write using fmt.Println either. So just ignore the error.
		_ = log.Sync()
	}()

	host := env.GetStringOrDefault("HOST", "127.0.0.1")
	port := env.GetIntOrDefault("PORT", 8080)

	s := server.NewServer(server.Options{
		Database: createDatabase(log),
		Host:     host,
		Port:     port,
		Log:      log,
		Queue:    createAzureServiceBus(log),
	})

	// create error group and context to run in background and listen for signals
	var errGroup errgroup.Group
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// following runs in background as a goroutine
	errGroup.Go(func() error {
		fmt.Println("context started listening for signals")
		<-ctx.Done()
		fmt.Println("context signal called.")
		if err := s.Stop(); err != nil {
			log.Info("Error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	// start server, also connect to DB and other stuff
	if err := s.Start(); err != nil {
		log.Info("Error starting server", zap.Error(err))
		return 1
	}

	// blocks here after server is started
	if err := errGroup.Wait(); err != nil {
		return 1
	}

	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func createAzureServiceBus(log *zap.Logger) *messaging.Queue {
	return messaging.NewQueue(
		messaging.NewQueueOptions{
			Namespace:   os.Getenv("AZURE_SERVICEBUS_NAMESPACE"),
			Log:         log,
			Name:        os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME"),
			KeyName:     os.Getenv("AZURE_SERVICEBUS_KEY_NAME"),
			KeyPassword: os.Getenv("AZURE_SERVICEBUS_KEY_VALUE"),
		},
	)
}

func createDatabase(log *zap.Logger) *storage.Database {
	return storage.NewDatabase(storage.NewDatabaseOptions{
		Host:                  env.GetStringOrDefault("DB_HOST", "localhost"),
		Port:                  env.GetIntOrDefault("DB_PORT", 5432),
		User:                  env.GetStringOrDefault("DB_USER", ""),
		Password:              env.GetStringOrDefault("DB_PASSWORD", ""),
		Name:                  env.GetStringOrDefault("DB_NAME", ""),
		MaxOpenConnections:    env.GetIntOrDefault("DB_MAX_OPEN_CONNECTIONS", 10),
		MaxIdleConnections:    env.GetIntOrDefault("DB_MAX_IDLE_CONNECTIONS", 10),
		ConnectionMaxLifetime: env.GetDurationOrDefault("DB_CONNECTION_MAX_LIFETIME", time.Hour),
		Log:                   log,
	})
}
