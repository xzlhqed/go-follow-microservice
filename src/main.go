package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"
	"github.com/xzlhqed/golang-follow-microservice/handlers"
)

func main() {

	l := log.New(os.Stdout, "user-api", log.LstdFlags)

	userHandler := handlers.ReturnUsers(l)

	serverMux := mux.NewRouter()

	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", userHandler.GetUsers)
	getRouter.HandleFunc("/{id:[0-9]+}", userHandler.GetUser)

	postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", userHandler.AddUser)
	postRouter.Use(userHandler.MiddlewareValidateUser)

	patchRouter := serverMux.Methods(http.MethodPatch).Subrouter()
	patchRouter.HandleFunc("/follow/{id1:[0-9]+}/{id2:[0-9]+}", userHandler.FollowUser)
	patchRouter.HandleFunc("/unfollow/{id1:[0-9]+}/{id2:[0-9]+}", userHandler.UnfollowUser)

	s := http.Server{
		Addr: ":9090",
		Handler: serverMux,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// listen for signal to close session
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}