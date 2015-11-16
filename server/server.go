package server

import (
	"fmt"
	templatepkg "github.com/defcube/webservice-proxy/server/templates"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
}

var templates *template.Template

func init() {
	var t *template.Template
	for _, fn := range templatepkg.AssetNames() {
		if strings.HasSuffix(fn, ".html") {
			if templates == nil {
				templates = template.New(fn)
				t = templates
			} else {
				t = templates.New(fn)
			}
			t.Parse(string(templatepkg.MustAsset(fn)))
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		s.handleProxy(w, r)
	} else if strings.HasPrefix(r.RequestURI, "/admin/") {
		templates.ExecuteTemplate(w, "index.html", nil)
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
