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
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
	"github.com/Carbon-X-DAO/QRInvite/templates"
	"github.com/mailgun/mailgun-go/v4"
)

var reInboundQR = regexp.MustCompile(`^\/qrcodes\/(?P<code>[0-9])$`)

const (
	queryInsertFormRow = `INSERT INTO
	form_info(
		first_name, last_name,
		country, department, city, neighborhood, street_address,
		id_no, phone, email, gender, age,
		daily_qty, weekly_qty, monthly_qty,
		newsletter, gift_box, authorized, claimed,
		ctime
	)
	VALUES (
		$1, $2,
		$3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12,
		$13, $14, $15,
		$16, $17, $18, $19,
		$20
	);`

	queryInsertQRIncomingHeaders = `INSERT INTO
	request_info(acceptlanguage, cookie, useragent, cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor, url_value, ctime)
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 )`

	queryInsertEmailStatus = `INSERT INTO
	email_status(email_address, gov_id, mailgun_msg, mailgun_id, error, ctime)
	VALUES( $1, $2, $3, $4, $5, $6)`
)

var stmtInsertQRIncomingHeaders *sql.Stmt
var stmtInsertFormRow *sql.Stmt
var stmtInsertEmailStatus *sql.Stmt

type Server struct {
	frontendRoot string
	*http.Server
	db    *sql.DB
	flyer image.Image
	mg    *mailgun.MailgunImpl
}

// tlsConfig may be nil, in which case an HTTP server will serve without TLS
func New(addr, mailgunAPIKey, flyerFilename, frontendRoot string, tlsConfig *tls.Config, db *sql.DB) (*Server, error) {
	var err error

	flyerHandle, err := os.Open(flyerFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open flyer image file: %s", err)
	}

	flyerImg, err := jpeg.Decode(flyerHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to JPEG decode flyer image file: %s", err)
	}

	mgClient := mailgun.NewMailgun("cieloverde.io", mailgunAPIKey)
	ds := mgClient.ListDomains(&mailgun.ListOptions{Limit: 20})

	var domains = []mailgun.Domain{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ds.Next(ctx, &domains)
	if err := ds.Err(); err != nil {
		return nil, fmt.Errorf("failed to use fresh mailgun client: %w", err)
	}

	server := &Server{
		frontendRoot: frontendRoot,
		db:           db,
		flyer:        flyerImg,
		mg:           mgClient,
	}

	if stmtInsertQRIncomingHeaders, err = db.Prepare(queryInsertQRIncomingHeaders); err != nil {
		return nil, fmt.Errorf("failed to prepare statement for storing incoming QR code handler headers: %w", err)
	}

	if stmtInsertFormRow, err = db.Prepare(queryInsertFormRow); err != nil {
		return nil, fmt.Errorf("failed to prepare statement for storing form information: %w", err)
	}

	if stmtInsertEmailStatus, err = db.Prepare(queryInsertEmailStatus); err != nil {
		return nil, fmt.Errorf("failed to prepare statement for storing email status information: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", server)

	httpServer := http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	server.Server = &httpServer

	return server, nil
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
	switch {
	case reInboundQR.MatchString(r.URL.Path) && r.Method == http.MethodGet:
		server.handleQRInbound(w, r)
	case r.URL.Path == "/submit" && r.Method == http.MethodPost:
		server.handleForm(w, r)
	default:
		server.handleFrontendPath(w, r)
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
		listing, err := fsutil.List(dir, server.frontendRoot)
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
