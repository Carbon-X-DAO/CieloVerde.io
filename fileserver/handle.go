package fileserver

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
)

func handleForm(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read /submit-form request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Printf("someone sent us form data!!! %s", bs)
	w.WriteHeader(http.StatusOK)
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
