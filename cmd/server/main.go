package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	//"os/signal"
	//"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

func main() {
	cfg := config.ServerConfigInit()

	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//определение роутера
	router := chi.NewRouter()
	//определение БД
	var store storage.Storage
	var err error
	//определение хендлера
	if len(cfg.DatabaseDSN) != 0 {
		client, err := postgresql.NewClient(context.Background(), 5, &cfg)
		checkError(err)
		store, err = repositories.NewPGSStore(client, &cfg)
		checkError(err)
	} else {
		store, err = repositories.NewRepository(&cfg)
		checkError(err)
	}

	handler := handlers.NewHandler(store)
	//регистрация хендлера
	handler.Register(router)

	srv := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	srv.ListenAndServe()

	//go func() {
	//	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()
	//log.Println("server is listen on", cfg.Addr)
	//
	////gracefullshutdown
	//<-done
	//
	//log.Print("Server Stopped")
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer func() {
	//	cancel()
	//}()
	//
	//if err = srv.Shutdown(ctx); err != nil {
	//	log.Fatalf("Server Shutdown Failed:%+v", err)
	//}
	//log.Print("Server Exited Properly")
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
