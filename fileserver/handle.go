package fileserver

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image/png"
	"io"
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
	queryInsertFormRow = `INSERT INTO
	form_info(
		first_name, last_name,
		country, department, city, town, neighborhood, street, street_number,
		id_no, phone, email, gender, age,
		daily_qty, weekly_qty, monthly_qty,
		newsletter, gift_box, authorized,
		ctime
	)
	VALUES (
		$1, $2,
		$3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12, $13, $14,
		$15,$16, $17,
		$18, $19, $20, 
		$21
	);`

	queryInsertQRIncomingHeaders = `INSERT INTO
	qr_incoming_headers(acceptlanguage, cookie, useragent, cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor, ctime)
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9 )`
)

var stmtInsertQRIncomingHeaders *sql.Stmt
var stmtInsertFormRow *sql.Stmt

type formInfo struct {
	FirstName    string `form:"fname"`
	LastName     string `form:"lname"`
	Country      string `form:"country"`
	Department   string `form:"department"`
	City         string `form:"city"`
	Town         string `form:"town"`
	Neighborhood string `form:"neighborhood"`
	Street       string `form:"street"`
	StreetNumber string `form:"address_number"`
	ID           uint64 `form:"id_no"`
	Phone        string `form:"phone"`
	Email        string `form:"email"`
	Gender       string `form:"gender"`
	Age          uint16 `form:"age"`
	DailyQty     string `form:"daily_qty"`
	WeeklyQty    string `form:"weekly_qty"`
	MonthlyQty   string `form:"monthly_qty"`
	Newsletter   bool   `form:"newsletter"`
	GiftBox      bool   `form:"gift_box"`
	Authorized   bool   `form:"authorized"`
}

func (server *Server) handleForm(w http.ResponseWriter, r *http.Request) {
	bs, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s\n", bs)
	r.Body = io.NopCloser(bytes.NewBuffer(bs))

	var fi formInfo
	dec := form.NewDecoder(r.Body)
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(&fi); err != nil {
		log.Printf("failed to decode form: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := saveFormInfo(ctx, &fi); err != nil {
		log.Printf("failed to save form: %+v: %s", fi, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("sending the email to %s\n", fi.Email)
	// TODO: validate email address
	go server.sendEmail(fi.Email, fi.ID)
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
	go saveQRInboundHeaders(r.Header)

	http.Redirect(w, r, "/form", http.StatusSeeOther)
}

func saveQRInboundHeaders(hdrs http.Header) {
	acceptlanguage := hdrs.Get("Accept-Language")
	cookie := hdrs.Get("Cookie")
	useragent := hdrs.Get("User-Agent")
	cfconnectingip := hdrs.Get("CF-Connecting-IP")
	xforwardedfor := hdrs.Get("X-Forwarded-For")
	cfray := hdrs.Get("CF-Ray")
	cfipcountry := hdrs.Get("CF-IPCountry")
	cfvisitor := hdrs.Get("CF-Visitor")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := stmtInsertQRIncomingHeaders.ExecContext(ctx, acceptlanguage, cookie, useragent, cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor, time.Now()); err != nil {
		log.Printf("failed to save inbound QR code headers: %s", err)
	}
}

func saveFormInfo(ctx context.Context, f *formInfo) error {
	_, err := stmtInsertFormRow.ExecContext(ctx, f.FirstName, f.LastName,
		f.Country, f.Department, f.City, f.Town, f.Neighborhood, f.Street, f.StreetNumber,
		f.ID, f.Phone, f.Email, f.Gender, f.Age,
		f.DailyQty, f.WeeklyQty, f.MonthlyQty,
		f.Newsletter, f.GiftBox, f.Authorized, time.Now())

	return err
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
