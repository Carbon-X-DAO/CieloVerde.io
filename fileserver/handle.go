package fileserver

import (
	"context"
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Carbon-X-DAO/CieloVerde.io/fsutil"
	"github.com/ajg/form"
)

type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

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

const tplNoSuchUser = `<!DOCTYPE html>
<html>
	<body>
		<h1>Usuario no existe</h1>
	</body>
</html>
`

const tplUnauthorized = `<!DOCTYPE html>
<html>
	<body>
		<h1>Desautorizado</h1>
	</body>
</html>
`

const loginForm = `
<!DOCTYPE html>
<html>
<style>
   input {
		font-size: 54px;
		display:block;
		margin: 40px;
   }
</style>
<body>
<form action="/login" method="POST">
<input name="username" placeholder="username" />
<input name="password" placeholder="password" />
<input type="submit" />
</form>
</body>
</html>
`

const tplAlreadyClaimed = `
<!DOCTYPE html>
<html>
	<body style="background-color: #F88685">
		<p> {{.First}} {{.Last}} </p>
		<p> {{.ID}} </p>
	</body>
</html>
`

const tplClaim = `
<!DOCTYPE html>
<html>
	<body style="background-color: #C8F5C6">
		<p> {{.First}} {{.Last}} </p>
		<p> {{.ID}} </p>
		<form  method="POST" action="/claim/{{.Hash}}">
		<input type="submit" value="reclamar" />
		</form>
	</body>
</html>
`

type user struct {
	First string
	Last  string
	ID    uint64
	Hash  string
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

	hash := md5.Sum([]byte(fmt.Sprintf("%d", fi.ID)))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := saveFormInfo(ctx, &fi, hash); err != nil {
		// has this ID already submitted an ID?
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Printf("failed to save form: %+v: %s", fi, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fi.Authorized {
		go saveRequestInfo(r.Header, r.URL)

		go func() {
			msg, id, err := server.sendEmail(fi.Email, hash)

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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (server *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	if _, err := w.Write([]byte(loginForm)); err != nil {
		log.Printf("failed to serve login form: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *Server) handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var li Login
	if err := form.NewDecoder(r.Body).Decode(&li); err != nil {
		log.Printf("failed to decode form from body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !(li.Username == server.adminUser && li.Password == server.adminPassword) {
		log.Println("rejected invalid login attempt")
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(tplUnauthorized))
		return
	}

	now := time.Now()
	until := now.Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "Shibboleth",
		Value:   server.shibboleth,
		Expires: until,
	})
}

func (server *Server) handleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(tplUnauthorized))
		return
	}

	for i, cook := range r.Cookies() {
		if cook.Name == "Shibboleth" && cook.Value == server.shibboleth {
			break
		}
		// we're on the last cookie and we didn't see what we need to see
		if i == len(r.Cookies())-1 {
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(tplUnauthorized))
			return
		}
	}

	hash := strings.TrimPrefix(r.URL.Path, "/users/")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rows, err := stmtSelectUser.QueryContext(ctx, hash)
	if err != nil {
		log.Printf("failed to select from form_info for hash %s: %s", hash, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var first string
	var last string
	var gov_id uint64
	var claimed bool

	var cnt int
	for rows.Next() {
		if err := rows.Scan(&first, &last, &gov_id, &claimed); err != nil {
			log.Printf("failed to scan user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cnt++
		if cnt > 1 {
			log.Printf("encountered multiple users with hash %s", hash)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if cnt == 0 {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(tplNoSuchUser))
		return
	}

	u := user{first, last, gov_id, hash}

	if claimed {
		t, err := template.New("alreadyClaimed").Parse(tplAlreadyClaimed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, u)
		return
	}

	t, err := template.New("claim").Parse(tplClaim)
	if err != nil {
		log.Printf("failed to generate template for claiming: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, u); err != nil {
		log.Printf("failed to execute template for claiming: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *Server) updateClaim(w http.ResponseWriter, r *http.Request) {
	hash := strings.TrimPrefix(r.URL.Path, "/claim/")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := stmtUpdateClaim.QueryContext(ctx, hash)
	if err != nil {
		log.Printf("failed to select from form_info for hash %s: %s", hash, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%s", hash), http.StatusSeeOther)
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

func saveFormInfo(ctx context.Context, f *formInfo, hash [16]byte) error {
	_, err := stmtInsertFormRow.ExecContext(ctx, f.FirstName, f.LastName,
		f.Country, f.Department, f.City, f.Neighborhood, f.Street,
		f.ID, f.Phone, f.Email, f.Gender, f.Age,
		f.DailyQty, f.WeeklyQty, f.MonthlyQty,
		f.Newsletter, f.GiftBox, f.Authorized, false, fmt.Sprintf("%x", hash), time.Now())

	return err
}
