package main

import (
	"context"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//// Relationship определяет положение в семье.
//type Relationship string
//
//// Возможные роли в семье.
//const (
//	Father      = Relationship("father")
//	Mother      = Relationship("mother")
//	Child       = Relationship("child")
//	GrandMother = Relationship("grandMother")
//	GrandFather = Relationship("grandFather")
//)
//
//// Family описывает семью.
//type Family struct {
//	Members map[Relationship]Person
//}
//
//// Person описывает конкретного человека в семье.
//type Person struct {
//	FirstName string
//	LastName  string
//	Age       int
//}
//
//var (
//	// ErrRelationshipAlreadyExists возвращает ошибку, если роль уже занята.
//	// Подробнее об ошибках поговорим в девятой теме: «Errors, log».
//	ErrRelationshipAlreadyExists = errors.New("relationship already exists")
//)
//
//// AddNew добавляет нового члена семьи.
//// Если в семье ещё нет людей, создаётся пустой map.
//// Если роль уже занята, метод выдаёт ошибку.
//func (f *Family) AddNew(r Relationship, p Person) error {
//	if f.Members == nil {
//		f.Members = map[Relationship]Person{}
//	}
//	if _, ok := f.Members[r]; ok {
//		return ErrRelationshipAlreadyExists
//	}
//	f.Members[r] = p
//	return nil
//}
//
//type User struct {
//	FirstName string
//	LastName  string
//}
//
//func (u User) FullName() string {
//	return u.FirstName + " " + u.LastName
//}
func main() {
	//f := Family{}
	//err := f.AddNew(Father, Person{
	//	FirstName: "Misha",
	//	LastName:  "Popov",
	//	Age:       56,
	//})
	//fmt.Println(f, err)
	//
	//err = f.AddNew(Father, Person{
	//	FirstName: "Drug",
	//	LastName:  "Mishi",
	//	Age:       57,
	//})
	//fmt.Println(f, err)
	//
	//v := Abs(3)
	//fmt.Println(v)
	//
	//u := User{
	//	FirstName: "Misha",
	//	LastName:  "Popov",
	//}
	//
	//fmt.Println(u.FullName())
	//
	//users := map[string]storage.User{
	//	"user1": {
	//		FirstName: "Test",
	//		LastName:  "Test",
	//	},
	//	"user2": {
	//		FirstName: "Test 2",
	//		LastName:  "Test 2",
	//	},
	//}
	//stor := repositories.NewRepository()
	//service := storage.NewService(stor)
	//http.HandleFunc("/user", handlers.UserViewHandler(users))
	//http.HandleFunc("/update/", handlers.UpdateHandler())
	//
	//log.Fatal(http.ListenAndServe(":8080", nil))
	serv := NewApp()
	serv.Run()
}

//func Abs(value float64) float64 {
//	return math.Abs(value)
//}

type App struct {
	server *http.Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {

	router := http.ServeMux{}
	repo := repositories.NewRepository()
	//service := storage.NewService(repo)
	handler := handlers.NewHandler(repo)
	handler.Register(&router)

	a.server = &http.Server{
		Addr:         ":8080",
		Handler:      &router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server is listen on port 8080")

	<-done

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
