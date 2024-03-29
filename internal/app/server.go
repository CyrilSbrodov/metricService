package app

import (
	"context"
	"crypto/rsa"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	"github.com/CyrilSbrodov/metricService.git/internal/app/protoserver"
	"github.com/CyrilSbrodov/metricService.git/internal/crypto"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

type ServerApp struct {
	router   *chi.Mux
	cfg      config.ServerConfig
	logger   *loggers.Logger
	cryptoer crypto.Cryptoer
	private  *rsa.PrivateKey
}

func NewServerApp() *ServerApp {
	cfg := config.ServerConfigInit()
	router := chi.NewRouter()
	logger := loggers.NewLogger()

	if cfg.CryptoPROKey != "" {
		c := crypto.NewCrypto()
		err := c.AddCryptoKey("public.pem", cfg.CryptoPROKey, "cert.pem", cfg.CryptoPROKeyPath)
		if err != nil {
			logger.LogErr(err, "filed to create file")
			os.Exit(1)
		}
		p, err := c.LoadPrivatePEMKey(cfg.CryptoPROKey)
		if err != nil {
			logger.LogErr(err, "filed to load file")
			os.Exit(1)
		}
		return &ServerApp{
			router:   router,
			cfg:      *cfg,
			logger:   logger,
			cryptoer: c,
			private:  p,
		}
	}
	return &ServerApp{
		router:   router,
		cfg:      *cfg,
		logger:   logger,
		cryptoer: nil,
		private:  nil,
	}
}

func (a *ServerApp) Run() {
	//определение БД
	store := protoserver.NewStorageServer(a.cfg, *a.logger)
	var err error
	//определение хендлера
	if len(a.cfg.DatabaseDSN) != 0 {
		client, err := postgresql.NewClient(context.Background(), 5, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
		store.Storage, err = repositories.NewPGSStore(client, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
	} else {
		store.Storage, err = repositories.NewRepository(&a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
	}
	if a.cfg.GRPCAddr != "" {
		store.Run()
	} else {
		handler := handlers.NewHandler(store.Storage, *a.logger, a.cryptoer, a.cfg, a.private)
		//регистрация хендлера
		handler.Register(a.router)

		srv := http.Server{
			Addr:    a.cfg.Addr,
			Handler: a.router,
		}
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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
}
