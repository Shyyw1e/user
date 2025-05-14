package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/Shyyw1e/user/internal/store"
	"github.com/Shyyw1e/user/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	st := store.NewInMemoryStore()
	h := handler.New(st)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 *time.Second))
	r.Use(middleware.Recoverer)

	r.Get("/users/", h.ListUsers)
	r.Post("/users/", h.CreateUser)
	r.Get("/users/{id}", h.GetUser)
	r.Put("/users/{id}", h.UpdateUser)
	r.Delete("/users/{id}", h.DeleteUser)

	srv := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Server Listening at localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err)

		}
	}()

	<-stop
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(),  5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}

	log.Println("Server exited properly")
}