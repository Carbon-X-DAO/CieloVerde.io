package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pg/pg/v10"
)

type App struct {
	db *pg.DB
}

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "john",
		Password: "",
		Database: "postgres",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to the postgres instance: %s", err)
	}
	cancel()

	a := App{
		db: db,
	}

	srv := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: a,
	}

	killed := make(chan os.Signal)
	signal.Notify(killed, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	serverShutdown := make(chan bool)

	go func() {
		sig := <-killed
		log.Printf("received signal to shutdown: %s", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown server: %s", err)
		}
		cancel()
		close(serverShutdown)
	}()

	log.Printf("starting the web server on port %d", 8080)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %s", err)
	}

	<-serverShutdown
	log.Printf("server has shut down... Exiting.")
}

func (a App) rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request to path %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.rootHandler(w, r)
}
