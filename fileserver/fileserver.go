package server

/*
Below code shamelessly taken from https://github.com/itsliamegan/fileserver
*/

import (
	"context"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Carbon-X-DAO/QRInvite/fsutil"
	"github.com/Carbon-X-DAO/QRInvite/templates"
)

type Server struct {
	root string
	mux  *http.ServeMux
	srv  *http.Server
}

func New(root string) *Server {
	mux := http.NewServeMux()
	server := &Server{root: root, mux: mux}
	mux.Handle("/", server)

	return server
}

func (server *Server) Listen(addr string) error {
	srv := http.Server{
		Addr:    addr,
		Handler: server.mux,
	}

	server.srv = &srv

	if err := server.srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.srv.Shutdown(ctx)
}

func (server *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/submit-form" && req.Method == http.MethodPost {
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("failed to read /submit-form request body: %s", err)
			res.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("someone sent us form data!!! %s", bs)
		res.WriteHeader(http.StatusOK)
		return
	}
	path := filepath.Join(server.root, req.URL.Path)
	exists, err := fsutil.Exists(path)
	if err != nil {
		writeErr(err, res)
		return
	}

	if exists {
		isDir, err := fsutil.IsDir(path)
		if err != nil {
			writeErr(err, res)
			return
		}

		if isDir {
			server.serveDir(path, res)
		} else {
			server.serveFile(path, res)
		}
	} else {
		htmlFile := path + ".html"
		exists, err = fsutil.Exists(htmlFile)
		if err != nil {
			writeErr(err, res)
			return
		}

		if exists {
			server.serveFile(htmlFile, res)
		} else {
			server.serveNotFound(res)
		}
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
