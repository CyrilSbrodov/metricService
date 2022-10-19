package main

import (
	"context"
	"flag"
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
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
)

var (
	flagAddress       *string
	flagRestore       *string
	flagStoreInterval *string
	flagStoreFile     *string
)

func init() {
	flagAddress = flag.String("a", "localhost:8080", "address of service")
	flagRestore = flag.String("r", "true", "restore from file")
	flagStoreInterval = flag.String("i", "300s", "upload interval")
	flagStoreFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file")
}

func main() {
	flag.Parse()
	cfg := config.NewConfigServer(*flagAddress, *flagStoreInterval, *flagStoreFile, *flagRestore)
	tickerUpload := time.NewTicker(cfg.StoreInterval)
	//определение роутера
	router := chi.NewRouter()
	//определение БД
	repo, err := repositories.NewRepository(cfg)
	if err != nil {
		fmt.Println(err)
	}
	//service := storage.NewService(repo)
	//определение хендлера
	handler := handlers.NewHandler(repo)
	//регистрация хендлера
	handler.Register(router)

	srv := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	//gracefullshutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//отправка данных на диск
	if cfg.StoreInterval != 0 {
		go uploadWithTicker(tickerUpload, repo)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server is listen on", cfg.Addr)

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

func uploadWithTicker(ticker *time.Ticker, repo *repositories.Repository) {
	for range ticker.C {
		err := repo.Upload()
		if err != nil {
			fmt.Println(err)
		}
	}
}
