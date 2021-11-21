package fileserver

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const (
	queryInsertFormRow = `INSERT INTO form_info(form, ctime) VALUES ($1, $2);`
)

type formInfo struct {
	Name    string
	Email   string
	Phone   string
	Message string
}

func (server *Server) handleCode(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/code/")
	hash := parts[1]

	code, err := qr.Encode(string(hash), qr.L, qr.Auto)
	if err != nil {
		log.Printf("failed to encode hash as QR code: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	intsize := 250
	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, intsize, intsize)
	if err != nil {
		log.Printf("failed to scale QR code: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, code); err != nil {
		log.Printf("failed to encode QR code as PNG: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Printf("failed to write HTTP response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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
		return
	}

	var form formInfo
	if err := json.NewDecoder(bytes.NewReader(bs)).Decode(&form); err != nil {
		log.Printf("failed to decode JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hash := md5.Sum([]byte(form.Email))

	log.Printf("redirecting to /code/%x : len(hash) == %d", hash, len(hash))

	w.Header().Set("Location", fmt.Sprintf("/code/%x", hash))

	w.WriteHeader(http.StatusSeeOther)

	if _, err := w.Write(nil); err != nil {
		log.Printf("failed to write HTTP response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("successfully returned response to form!!!")
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
