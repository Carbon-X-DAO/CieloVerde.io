package fileserver

/*
Below code shamelessly taken from https://github.com/itsliamegan/fileserver
*/

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
	"github.com/Carbon-X-DAO/QRInvite/templates"
)

type Server struct {
	root string
	*http.Server
	db *sql.DB
	// tlsConfig *tls.Config
}

// tlsConfig may be nil, in which case an HTTP server will serve without TLS
func New(addr string, root string, tlsConfig *tls.Config, db *sql.DB) *Server {
	server := &Server{root: root, db: db}

	mux := http.NewServeMux()
	mux.Handle("/", server)

	httpServer := http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	server.Server = &httpServer

	return server
}

func (server *Server) Listen() error {
	if server.Server.TLSConfig != nil {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("TLS HTTP server failed: %s", err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTP server failed: %s", err)
		}
	}

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	if err := server.Server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to shut shut down HTTP server: %s", err)
	}
	if err := server.db.Close(); err != nil {
		return fmt.Errorf("failed to close DB connection: %s", err)
	}

	return nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling a request at path: %s", r.URL.Path)

	isCodePath, err := regexp.MatchString(`^\/code\/[a-f0-9]{32}$`, r.URL.Path)
	if err != nil {
		log.Printf("failed to compare path to regex: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch {
	case r.URL.Path == "/submit" && r.Method == http.MethodPost:
		server.handleForm(w, r)
	case isCodePath && r.Method == http.MethodGet:
		server.handleCode(w, r)
	default:
		server.handlePath(w, r)
	}
}

func (server *Server) serveFile(file string, res http.ResponseWriter) {
	fp, err := os.Open(file)
	if err != nil {
		writeErr(err, res)
		return
	}
	defer fp.Close()

	_, err = io.Copy(res, fp)
	if err != nil {
		writeErr(err, res)
	}
}

func (server *Server) serveDir(dir string, res http.ResponseWriter) {
	indexFile := filepath.Join(dir, "index.html")
	exists, err := fsutil.Exists(indexFile)
	if err != nil {
		writeErr(err, res)
	}

	if exists {
		server.serveFile(indexFile, res)
	} else {
		listing, err := fsutil.List(dir, server.root)
		if err != nil {
			writeErr(err, res)
			return
		}

		writeTemplate(templates.Listing, listing, res)
	}
}

func (server *Server) serveNotFound(res http.ResponseWriter) {
	res.WriteHeader(http.StatusNotFound)
	writeTemplate(templates.NotFound, nil, res)
}

func writeTemplate(tmpl *template.Template, ctx interface{}, res http.ResponseWriter) {
	err := tmpl.Execute(res, ctx)
	if err != nil {
		writeErr(err, res)
	}
}

func writeErr(err error, res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)
	templates.Error.Execute(res, err)
	log.Printf("err: %s", err.Error())
}
