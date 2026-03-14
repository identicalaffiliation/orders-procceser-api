package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/broker"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/config"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/logger"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/api"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/storage"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/service"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/transport/rest"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg := config.MustLoadConfig(configPath)

	slogger := logger.NewLogger(cfg)

	psql, err := storage.NewConnect(cfg)
	if err != nil {
		slogger.Error("connect to postgres", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := psql.Close(); err != nil {
			slogger.Error("close postgres", "error", err)
		}
	}()

	apiRepository := api.NewAPIRepository(psql, slogger)

	publisher, err := broker.NewBroker(cfg)
	if err != nil {
		slogger.Error("init rabbitmq", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := publisher.Close(); err != nil {
			slogger.Error("close rabbitmq", "error", err)
		}
	}()

	service := service.NewOrderService(apiRepository, slogger, publisher)
	server := rest.NewOrdersAPI(service, slogger, cfg)
	server.SetupAPI()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.StartAPI(); err != nil {
			slogger.Error("start api", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := server.ShutdownAPI(shutdownCtx); err != nil {
		slogger.Error("shutdown api", "error", err)
		os.Exit(1)
	}

	slogger.Debug("server was stopped gracefully")
}
