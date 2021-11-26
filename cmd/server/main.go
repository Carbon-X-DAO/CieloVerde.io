package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	fileserver "github.com/Carbon-X-DAO/QRInvite/fileserver"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	flagAddress  string
	flagRoot     string
	flagDBRole   string
	flagCertFile string
	flagKeyFile  string
)

func init() {
	flag.StringVar(&flagAddress, "address", "0.0.0.0:80", "address on which to listen")
	flag.StringVar(&flagRoot, "root", "/result/static", "root path to site")
	flag.StringVar(&flagDBRole, "role", "postgres", "postgres DB user role")
	flag.StringVar(&flagCertFile, "cert", "example.crt", "TLS certificate file")
	flag.StringVar(&flagKeyFile, "key", "example.key", "TLS certificate signing key file")
	flag.Parse()
}

func main() {
	dbHost := "localhost"
	dbPort := 5432
	dbName := "postgres"
	opts := "sslmode=disable"

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s@%s:%d/%s?%s", flagDBRole, dbHost, dbPort, dbName, opts))
	if err != nil {
		log.Fatalf("failed to initialize a postgres instance for DBMS configuration: %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to the postgres instance: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := db.ExecContext(ctx, `CREATE DATABASE qrinvite`); err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatalf("failed to create database qrinvite: %s", err)
	}
	cancel()

	// switch to the qrinvite DB now that it has been created
	dbName = "qrinvite"
	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s@%s:%d/%s?%s", flagDBRole, dbHost, dbPort, dbName, opts))
	if err != nil {
		log.Fatalf("failed to initialize a postgres instance for app usage: %s", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: "qrinvite",
	})
	if err != nil {
		log.Fatalf("failed to obtain postgres driver for migrations: %s", err)
	}

	// apply migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"qrinvite",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to initialize a migrate driver instance: %s", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply all migrations: %s", err)
	}

	var root string

	if len(flag.Args()) > 0 {
		root = flag.Arg(0)
	} else {
		root, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	root += flagRoot
	log.Printf("starting a fileserver for root path: %s", root)

	cert, err := tls.LoadX509KeyPair(flagCertFile, flagKeyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s and %s: %s", flagCertFile, flagKeyFile, err)
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	srv := fileserver.New(root, &tlsConfig, db)

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

	log.Printf("starting the web server on address %s", flagAddress)
	if err := srv.Listen(flagAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %s", err)
	}

	<-serverShutdown
	log.Printf("server has shut down... Exiting.")
}
