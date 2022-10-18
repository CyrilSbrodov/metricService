package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
)

func main() {

	//определение роутера
	router := chi.NewRouter()
	//определение БД
	repo := repositories.NewRepository()
	//service := storage.NewService(repo)
	//определение хендлера
	handler := handlers.NewHandler(repo)
	//регистрация хендлера
	handler.Register(router)

	//srv := http.Server{
	//	Addr:    ":8080",
	//	Handler: router,
	//}
	if err := http.ListenAndServe(":8080", router); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	//gracefullshutdown
	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//go func() {
	//	if err := srv.ListenAndServe(":8080", router); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()
	//log.Println("server is listen on port 8080")
	//
	//<-done
	//
	//log.Print("Server Stopped")
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer func() {
	//	cancel()
	//}()

	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatalf("Server Shutdown Failed:%+v", err)
	//}
	//log.Print("Server Exited Properly")
}
