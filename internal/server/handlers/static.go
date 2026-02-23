package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/jsnfwlr/go11y"
)

// UI inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h Handlers) UI(w http.ResponseWriter, r *http.Request) {
	_, o := go11y.Get(r.Context())

	// Join internally call path.Clean to prevent directory traversal

	if r.URL.Path == "/" {
		http.ServeFileFS(w, r, h.staticFS, "index.html")
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/")

	// check whether a file exists - if not serve the 404.html, if we have any other error
	// return a 500 internal server error
	_, err := h.staticFS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.ServeFileFS(w, r, h.staticFS, "index.html")
			return
		}
		o.Error("error with embedded file system", err, go11y.SeverityHigh, "file_path", path)

		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// otherwise, use http.FileServer to serve the static file
	// http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
	http.ServeFileFS(w, r, h.staticFS, path)
}
