package server

import (
	"fmt"
	"github.com/defcube/webservice-proxy/server/static"
	templatepkg "github.com/defcube/webservice-proxy/server/templates"
	"github.com/elazarl/go-bindata-assetfs"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

// todo remove all bindata

type Server struct {
	syncInitOnce sync.Once

	// staticServer serves from /static/
	staticServer http.Handler

	// templates contains all the imported templates from /templates/
	templates *template.Template
}

func (s *Server) Init() {
	s.syncInitOnce.Do(func() {
		s.staticServer = http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, Prefix: ""})
		s.staticServer = http.StripPrefix("/static", s.staticServer)

		var t *template.Template
		for _, fileName := range templatepkg.AssetNames() {
			if strings.HasSuffix(fileName, ".html") {
				if s.templates == nil {
					s.templates = template.New(fileName)
					s.templates.Funcs(map[string]interface{}{
						"safeHTML": func(s interface{}) template.HTML {
							return template.HTML(fmt.Sprintf("%v", s))
						},
					})
					t = s.templates
				} else {
					t = s.templates.New(fileName)
				}
				template.Must(t.Parse(string(templatepkg.MustAsset(fileName))))
			}
		}
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Init()
	if r.RequestURI == "/" {
		s.handleProxy(w, r)
	} else if r.RequestURI == "/admin/form/" {
		s.handleFormTest(w, r)
	} else if r.RequestURI == "/admin/" {
		s.handleAdmin(w, r)
	} else if strings.HasPrefix(r.RequestURI, "/static/") {
		s.staticServer.ServeHTTP(w, r)
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 not found")
	}
}
