package main

import (
	parserpkg "apiDocTransfer/internal/service/parser"
	"apiDocTransfer/internal/service/postman_client"
	"apiDocTransfer/pkg/config"
	pkgLogger "apiDocTransfer/pkg/logger"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	appConfig, err := config.Load()
	if err != nil {
		log.Fatal("error loading config:", err)
	}

	logger, loggerDefer := initLogger(appConfig.LogLevel)
	defer loggerDefer(logger)

	logger.Debug("config", zap.Any("config", appConfig))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	postmanHttpClient := &http.Client{
		Timeout: appConfig.PostmanApiConnectionTimeout,
	}

	postmanClient := postman_client.NewPostmanClient(
		appConfig.PostmanApiKey,
		appConfig.PostmanWorkspaceID,
		postmanHttpClient,
	)

	parser := parserpkg.NewParser(*appConfig, logger, client, postmanClient)
	if err := parser.Parse(ctx); err != nil {
		logger.Error("parse failed", zap.Error(err))
		return
	}
	logger.Info("parsed successfully")
}

func initLogger(logLevel string) (logger *zap.Logger, loggerDefer func(logger *zap.Logger)) {
	loggerOptions := pkgLogger.Options{
		Level: logLevel,
	}

	loggerDefer = func(logger *zap.Logger) {
		func(logger *zap.Logger) {
			err := logger.Sync()
			if err != nil && !errors.Is(err, syscall.EINVAL) {
				log.Fatalf("failed to sync logger: %v", err)
			}
		}(logger)
	}

	return pkgLogger.New(loggerOptions), loggerDefer
}
