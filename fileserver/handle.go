package fileserver

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
	"github.com/ajg/form"
)

const (
	queryInsertFormRow = `INSERT INTO
	form_info(
		first_name, last_name,
		country, department, city, town, neighborhood, street, street_number,
		id_no, phone, email, gender, age,
		daily_qty, weekly_qty, monthly_qty,
		newsletter, gift_box, authorized, claimed,
		ctime
	)
	VALUES (
		$1, $2,
		$3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12, $13, $14,
		$15,$16, $17,
		$18, $19, $20, $21,
		$22
	);`

	queryInsertQRIncomingHeaders = `INSERT INTO
	request_info(acceptlanguage, cookie, useragent, cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor, url_value, ctime)
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 )`
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
		if strings.Contains(err.Error(), "duplicate") {
			return
		}
		log.Printf("failed to save form: %+v: %s", fi, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go saveRequestInfo(r.Header, r.URL)

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
	go saveRequestInfo(r.Header, r.URL)

	http.Redirect(w, r, "/form", http.StatusSeeOther)
}

func saveRequestInfo(hdrs http.Header, url *url.URL) {
	acceptlanguage := hdrs.Get("Accept-Language")
	cookie := hdrs.Get("Cookie")
	useragent := hdrs.Get("User-Agent")
	cfconnectingip := hdrs.Get("CF-Connecting-IP")
	xforwardedfor := hdrs.Get("X-Forwarded-For")
	cfray := hdrs.Get("CF-Ray")
	cfipcountry := hdrs.Get("CF-IPCountry")
	cfvisitor := hdrs.Get("CF-Visitor")
	u := url.String()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := stmtInsertQRIncomingHeaders.ExecContext(ctx,
		acceptlanguage, cookie, useragent,
		cfconnectingip, xforwardedfor, cfray, cfipcountry, cfvisitor,
		u,
		time.Now()); err != nil {
		log.Printf("failed to save request infos: %s", err)
	}
}

func saveFormInfo(ctx context.Context, f *formInfo) error {
	_, err := stmtInsertFormRow.ExecContext(ctx, f.FirstName, f.LastName,
		f.Country, f.Department, f.City, f.Town, f.Neighborhood, f.Street, f.StreetNumber,
		f.ID, f.Phone, f.Email, f.Gender, f.Age,
		f.DailyQty, f.WeeklyQty, f.MonthlyQty,
		f.Newsletter, f.GiftBox, f.Authorized, false, time.Now())

	return err
}
