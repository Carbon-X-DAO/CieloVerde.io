package fileserver

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
	"github.com/ajg/form"
)

type formInfo struct {
	FirstName    string `form:"fname"`
	LastName     string `form:"lname"`
	Country      string `form:"country"`
	Department   string `form:"department"`
	City         string `form:"city"`
	Neighborhood string `form:"neighborhood"`
	Street       string `form:"street_address"`
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
		// has this ID already submitted an ID?
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Printf("failed to save form: %+v: %s", fi, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go saveRequestInfo(r.Header, r.URL)

	go func() {
		msg, id, err := server.sendEmail(fi.Email, fi.ID)

		var errString string
		if err != nil {
			errString = err.Error()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if _, err := stmtInsertEmailStatus.ExecContext(ctx, fi.Email, fi.ID, msg, id, errString, time.Now()); err != nil {
			log.Printf("failed to store email status info in DB (%s, %d): %s", fi.Email, fi.ID, err)
		}
	}()
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
		f.Country, f.Department, f.City, f.Neighborhood, f.Street,
		f.ID, f.Phone, f.Email, f.Gender, f.Age,
		f.DailyQty, f.WeeklyQty, f.MonthlyQty,
		f.Newsletter, f.GiftBox, f.Authorized, false, time.Now())

	return err
}
