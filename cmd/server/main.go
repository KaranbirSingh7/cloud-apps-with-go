// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"canvas/server"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
	log, err := createLogger(getStringOrDefault("LOG_ENV", "development"))
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

	host := getStringOrDefault("HOST", "127.0.0.1")
	port := getIntOrDefault("PORT", 8080)

	s := server.NewServer(server.Options{
		Host: host,
		Port: port,
		Log:  log,
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

	// start server
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

func getStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}

func getIntOrDefault(name string, defaultV int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	vAsInt, err := strconv.Atoi(v)
	if err != nil {
		return defaultV
	}
	return vAsInt
}
