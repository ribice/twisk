package openapi

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	_ "github.com/ribice/twisk/openapi/statik" // Static files
)

// New creates new openapi http service
func New(r *mux.Router, uname, password string) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	svc := Openapi{
		username:     uname,
		password:     password,
		staticServer: http.FileServer(statikFS),
	}

	r.PathPrefix("/openapi").HandlerFunc(svc.docHandler)
}

// Openapi represents openapi documentation service
type Openapi struct {
	username     string
	password     string
	staticServer http.Handler
}

func (o *Openapi) docHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimRight(r.URL.Path, "/")
	if path == "/openapi/swaggerui" {
		r.URL.Path = "/openapi/swagger.html"
	}
	w.Header().Set("Content-Type", o.getContentType(path))

	http.StripPrefix("/openapi/", o.serve()).ServeHTTP(w, r)
}

func (o *Openapi) serve() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || (u != o.username || p != o.password) {
			w.Header().Add("WWW-Authenticate", `Basic realm="Access to swagger docs"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		o.staticServer.ServeHTTP(w, r)
	})
}

// Workaround for bug with wrong mime type returned
func (o *Openapi) getContentType(path string) string {
	split := strings.Split(path, ".")
	if len(split) < 2 {
		return "text/html"
	}
	switch split[1] {
	case "css":
		return "text/css"
	case "png":
		return "image/png"
	default:
		return "text/html"
	}
}
