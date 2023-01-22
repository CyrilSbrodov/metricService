package app

import (
	"context"
	"net/http"
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

type ServerApp struct {
	router *chi.Mux
	cfg    config.ServerConfig
	logger *loggers.Logger
}

func NewServerApp() *ServerApp {
	router := chi.NewRouter()
	logger := loggers.NewLogger()
	cfg := config.ServerConfigInit()
	return &ServerApp{
		router: router,
		cfg:    cfg,
		logger: logger,
	}
}

func (a *ServerApp) Run() {
	//определение БД
	var store storage.Storage
	var err error
	//определение хендлера
	if len(a.cfg.DatabaseDSN) != 0 {
		client, err := postgresql.NewClient(context.Background(), 5, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
		store, err = repositories.NewPGSStore(client, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
	} else {
		store, err = repositories.NewRepository(&a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
	}

	handler := handlers.NewHandler(store, a.logger)
	//регистрация хендлера
	handler.Register(a.router)

	srv := http.Server{
		Addr:    a.cfg.Addr,
		Handler: a.router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.LogErr(err, "server not started")
		}
	}()
	a.logger.LogInfo("server is listen:", a.cfg.Addr, "start server")

	//gracefullshutdown
	<-done

	a.logger.LogInfo("", "", "server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctx); err != nil {
		a.logger.LogErr(err, "Server Shutdown Failed")
	}
	a.logger.LogInfo("", "", "Server Exited Properly")
}
