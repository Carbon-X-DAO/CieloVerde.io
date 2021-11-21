package fileserver

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
)

const (
	queryInsertFormRow = `INSERT INTO form_info(form, ctime) VALUES ($1, $2);`
)

func (server *Server) handleForm(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read /submit-form request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("someone sent us form data!!! %s", bs)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err = server.db.ExecContext(ctx, queryInsertFormRow, bs, time.Now()); err != nil {
		log.Printf("failed to read /submit-form request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) handlePath(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(server.root, r.URL.Path)
	exists, err := fsutil.Exists(path)
	if err != nil {
		writeErr(err, w)
		return
	}

	if exists {
		isDir, err := fsutil.IsDir(path)
		if err != nil {
			writeErr(err, w)
			return
		}

		if isDir {
			server.serveDir(path, w)
		} else {
			server.serveFile(path, w)
		}
	} else {
		htmlFile := path + ".html"
		exists, err = fsutil.Exists(htmlFile)
		if err != nil {
			writeErr(err, w)
			return
		}

		if exists {
			server.serveFile(htmlFile, w)
		} else {
			server.serveNotFound(w)
		}
	}
}
