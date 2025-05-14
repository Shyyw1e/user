package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Shyyw1e/user/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handler.ListUsers)
		r.Post("/", handler.CreateUser)
		r.Get("/{id}", handler.GetUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})

	srv := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server Listening at localhost:8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listening: %s\n", err)

		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(),  5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}

	log.Println("Server exited properly")
}