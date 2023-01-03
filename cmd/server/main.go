package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

func main() {
	cfg := config.ServerConfigInit()
	logger := loggers.NewLogger()

	//определение роутера
	router := chi.NewRouter()
	//определение БД
	var store storage.Storage
	var err error
	//определение хендлера
	if len(cfg.DatabaseDSN) != 0 {
		client, err := postgresql.NewClient(context.Background(), 5, &cfg, logger)
		if err != nil {
			logger.LogErr(err, "")
			os.Exit(1)
		}
		store, err = repositories.NewPGSStore(client, &cfg, logger)
		if err != nil {
			logger.LogErr(err, "")
			os.Exit(1)
		}
	} else {
		store, err = repositories.NewRepository(&cfg, logger)
		if err != nil {
			logger.LogErr(err, "")
			os.Exit(1)
		}
	}

	handler := handlers.NewHandler(store, logger)
	//регистрация хендлера
	handler.Register(router)

	srv := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogErr(err, "server not started")
		}
	}()
	logger.LogInfo("server is listen:", cfg.Addr, "start server")

	//gracefullshutdown
	<-done

	logger.LogInfo("", "", "server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctx); err != nil {
		logger.LogErr(err, "Server Shutdown Failed")
	}
	logger.LogInfo("", "", "Server Exited Properly")
}
