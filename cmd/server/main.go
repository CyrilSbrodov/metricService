package main

import (
	"context"
	//"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

func main() {
	cfg := config.ServerConfigInit()
	//tickerUpload := time.NewTicker(cfg.StoreInterval)
	fmt.Println(cfg.DatabaseDSN)
	//определение роутера
	router := chi.NewRouter()
	//определение БД
	var store storage.Storage
	client, err := postgresql.NewClient(context.Background(), 5, &cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//определение хендлера
	if cfg.DatabaseDSN != "" {
		store, err = repositories.NewPGSStore(client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		store, err = repositories.NewRepository(&cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	handler := handlers.NewHandler(store)
	//регистрация хендлера
	handler.Register(router)

	srv := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	//gracefullshutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server is listen on", cfg.Addr)

	//отправка данных на диск, если запись разрешена и файл создан
	//if store.Check {
	//	go uploadWithTicker(tickerUpload, repo, done)
	//}

	<-done

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func uploadWithTicker(ticker *time.Ticker, repo *repositories.Repository, done chan os.Signal) {
	for {
		select {
		case <-ticker.C:
			err := repo.Upload()
			if err != nil {
				fmt.Println(err)
				return
			}
		case <-done:
			ticker.Stop()
			return
		}
	}
}
