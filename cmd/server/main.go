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

	fileserver "github.com/Carbon-X-DAO/CieloVerde.io/fileserver"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	flagAddress        string
	flagAdminUser      string
	flagAdminPassword  string
	flagMailgunAPIKey  string
	flagFlyerFilename  string
	flagDBName         string
	flagRoot           string
	flagDBRole         string
	flagCertFile       string
	flagKeyFile        string
	flagShibbolethGUID string
)

func init() {
	flag.StringVar(&flagAddress, "address", "0.0.0.0:80", "address on which to listen")
	flag.StringVar(&flagAdminUser, "adminuser", "cielo", "username of admin")
	flag.StringVar(&flagAdminPassword, "adminpassword", "verde", "password of admin")

	flag.StringVar(&flagMailgunAPIKey, "mg", "", "priavte Mailgun API key")
	flag.StringVar(&flagRoot, "root", "./result/static", "root path to site")
	flag.StringVar(&flagFlyerFilename, "flyer", "./flyer.jpg", "path to flyer image")
	flag.StringVar(&flagDBName, "dbname", "", "name of DB")
	flag.StringVar(&flagDBRole, "role", "postgres", "postgres DB user role")
	flag.StringVar(&flagCertFile, "cert", "example.crt", "TLS certificate file")
	flag.StringVar(&flagKeyFile, "key", "example.key", "TLS certificate signing key file")
	flag.StringVar(&flagShibbolethGUID, "guid", "5f9f3021-845e-4a21-a63c-2f7f75262649", "a special token that admins need")
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
	if _, err := db.ExecContext(ctx, fmt.Sprintf(`CREATE DATABASE %s`, flagDBName)); err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatalf("failed to create database %s: %s", flagDBName, err)
	}
	cancel()

	// switch to the %s DB now that it has been created
	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s@%s:%d/%s?%s", flagDBRole, dbHost, dbPort, flagDBName, opts))
	if err != nil {
		log.Fatalf("failed to initialize a postgres instance for app usage: %s", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: flagDBName,
	})
	if err != nil {
		log.Fatalf("failed to obtain postgres driver for migrations: %s", err)
	}

	// apply migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		flagDBName,
		driver,
	)
	if err != nil {
		log.Fatalf("failed to initialize a migrate driver instance: %s", err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("failed to apply all migrations: %s", err)
		}
		log.Println("DB already up to date :)")
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

	var tlsConfig *tls.Config
	if flagCertFile != "" && flagKeyFile != "" {
		log.Printf("cert file is empty: %v", flagCertFile == "")
		log.Printf("key file is empty: %v", flagKeyFile == "")
		cert, err := tls.LoadX509KeyPair(flagCertFile, flagKeyFile)
		if err != nil {
			log.Fatalf("failed to load key pair: %s and %s: %s", flagCertFile, flagKeyFile, err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	// XOR trick from https://stackoverflow.com/a/23025720
	if (flagCertFile != "") != (flagKeyFile != "") {
		log.Fatal("both cert file and key file must be either non-empty or empty")
	}

	srv, err := fileserver.New(flagAddress, flagAdminUser, flagAdminPassword, flagShibbolethGUID, flagMailgunAPIKey, flagFlyerFilename, root, tlsConfig, db)
	if err != nil {
		log.Fatalf("failed to create a new fileserver instance: %s", err)
	}

	log.Printf("starting a fileserver for root path: %s", root)

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
	if err := srv.Listen(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %s", err)
	}

	<-serverShutdown
	log.Printf("server has shut down... Exiting.")
}
