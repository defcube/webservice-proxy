package server

import (
	"fmt"
	"github.com/bluele/gforms"
	"github.com/defcube/webservice-proxy/server/static"
	templatepkg "github.com/defcube/webservice-proxy/server/templates"
	"github.com/elazarl/go-bindata-assetfs"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

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
		for _, fn := range templatepkg.AssetNames() {
			if strings.HasSuffix(fn, ".html") {
				if s.templates == nil {
					s.templates = template.New(fn)
					s.templates.Funcs(map[string]interface{}{
						"safeHTML": func(s interface{}) template.HTML {
							return template.HTML(fmt.Sprintf("%v", s))
						},
					})
					t = s.templates
				} else {
					t = s.templates.New(fn)
				}
				template.Must(t.Parse(string(templatepkg.MustAsset(fn))))
			}
		}
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Init()
	if r.RequestURI == "/" {
		s.handleProxy(w, r)
	} else if strings.HasPrefix(r.RequestURI, "/admin/") {
		s.handleAdmin(w, r)
	} else if strings.HasPrefix(r.RequestURI, "/static/") {
		s.staticServer.ServeHTTP(w, r)
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 not found")
	}
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.PostForm
	url := form.Get("-url")
	if url == "" {
		panic("missing url") // todo handle gracefully
	}
	form.Del("-url")
	client := http.Client{}
	resp, err := client.PostForm(url, form)
	if err != nil {
		panic(err) // todo handle
	}
	w.WriteHeader(resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) // todo
	}
	w.Write(respBody)
}

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	f := gforms.DefineForm(gforms.NewFields(
		gforms.NewTextField("Foo", nil, nil),
	))(r)
	fi := f.Fields()[0]
	err := s.templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Form": f, "Field": fi, "H": template.HTML(fi.Html())})
	if err != nil {
		panic(err)
	}
}
