package fileserver

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
	"github.com/ajg/form"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const (
	queryInsertFormRow = `INSERT INTO form_info(form, ctime) VALUES ($1, $2);`

	queryInsertQRIncomingHeaders = "INSERT INTO qr_incoming_headers(acceptlanguage, cookie, useragent, cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor, ctime) VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9 )"
)

var stmtInsertQRIncomingHeaders *sql.Stmt
var stmtInsertFormRow *sql.Stmt

type formInfo struct {
	Name    string
	Email   string
	Phone   string
	Message string
}

func (server *Server) handleForm(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read /submit request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err = stmtInsertFormRow.ExecContext(ctx, bs, time.Now()); err != nil {
		log.Printf("failed to read /submit request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var fi formInfo
	dec := form.NewDecoder(bytes.NewReader(bs))
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(&fi); err != nil {
		log.Printf("failed to decode form: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hash := md5.Sum([]byte(fi.Email))
	hashString := string(fmt.Sprintf("%x", hash))

	go saveQRCodePNG(server.ticketsDir, hashString)

	if _, err := w.Write([]byte(fmt.Sprintf(ticketHTMLTemplate, ticketFilename(server.ticketsDir, hashString)))); err != nil {
		log.Printf("failed to execute ticket template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *Server) handleFrontendPath(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(server.frontendRoot, r.URL.Path)
	exists, err := fsutil.Exists(path)
	if err != nil {
		writeErr(err, w)
		return
	}

	w.Header().Add("Cache-Control", "max-age=86400,s-maxage=86400")
	if exists {
		isDir, err := fsutil.IsDir(path)
		if err != nil {
			writeErr(err, w)
			return
		}

		if isDir {
			server.serveDir(path, w)
			return
		}

		switch filepath.Ext(path) {
		case ".css":
			{
				w.Header().Set("Content-Type", "text/css")
			}
		case ".js":
			{
				w.Header().Set("Content-Type", "application/javascript")
			}
		}

		server.serveFile(path, w)
		return
	}

	htmlFile := path + ".html"
	exists, err = fsutil.Exists(htmlFile)
	if err != nil {
		writeErr(err, w)
		return
	}

	if exists {
		w.Header().Set("Content-Type", "text/html")
		server.serveFile(htmlFile, w)
	} else {
		server.serveNotFound(w)
	}
}

func (server *Server) handleTicketsPath(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("./", r.URL.Path)
	exists, err := fsutil.Exists(path)
	if err != nil {
		writeErr(err, w)
		return
	}

	w.Header().Add("Cache-Control", "max-age=86400,s-maxage=86400")
	if exists {
		isDir, err := fsutil.IsDir(path)
		if err != nil {
			writeErr(err, w)
			return
		}

		if isDir {
			server.serveDir(path, w)
			return
		}

		switch filepath.Ext(path) {
		case ".css":
			{
				w.Header().Set("Content-Type", "text/css")
			}
		case ".js":
			{
				w.Header().Set("Content-Type", "application/javascript")
			}
		}

		server.serveFile(path, w)
		return
	}

	htmlFile := path + ".html"
	exists, err = fsutil.Exists(htmlFile)
	if err != nil {
		writeErr(err, w)
		return
	}

	if exists {
		w.Header().Set("Content-Type", "text/html")
		server.serveFile(htmlFile, w)
	} else {
		server.serveNotFound(w)
	}
}

func (server *Server) handleQRInbound(w http.ResponseWriter, r *http.Request) {
	go saveQRInboundHeaders(r)

	http.Redirect(w, r, "/form", http.StatusSeeOther)
}

func saveQRInboundHeaders(r *http.Request) {
	acceptlanguage := r.Header.Get("Accept-Language")
	cookie := r.Header.Get("Cookie")
	useragent := r.Header.Get("User-Agent")
	cfconnectingip := r.Header.Get("CF-Connecting-IP")
	xforwardedfor := r.Header.Get("X-Forwarded-For")
	cfray := r.Header.Get("CF-Ray")
	cfipcountry := r.Header.Get("CF-IPCountry")
	cfvisitor := r.Header.Get("CF-Visitor")

	// queryInsertQRIncomingHeaders = "INSERT INTO qr_incoming_headers(acceptlanguage, cookie, useragent, cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor, ctime) VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9 )"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := stmtInsertQRIncomingHeaders.ExecContext(ctx, acceptlanguage, cookie, useragent, cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor, time.Now()); err != nil {
		log.Printf("failed to save inbound QR code headers: %s", err)
	}
}

func saveQRCodePNG(dir, hash string) {
	code, err := qr.Encode(string(hash), qr.L, qr.Auto)
	if err != nil {
		log.Printf("failed to encode hash as QR code: %s\n", err)
		return
	}

	intsize := 256
	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, intsize, intsize)
	if err != nil {
		log.Printf("failed to scale QR code: %s\n", err)
		return
	}

	filelocation := ticketFilename(dir, hash)
	f, err := os.Create(filelocation)
	if err != nil {
		log.Printf("failed to create file %s: %s", filelocation, err)
		return
	}

	if err := png.Encode(f, code); err != nil {
		log.Printf("failed to write QR code as PNG to file %s: %s", filelocation, err)
		return
	}
}

func ticketFilename(dir, name string) string {
	filename := fmt.Sprintf("%s.png", name)
	return filepath.Join(dir, filename)
}
