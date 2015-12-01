package server

import (
	"fmt"
	"github.com/defcube/webservice-proxy/server/internal/static"
	"github.com/defcube/webservice-proxy/server/internal/stats"
	templatepkg "github.com/defcube/webservice-proxy/server/internal/templates"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/namsral/flag"
	"github.com/simonz05/stathat"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	stathatKey = flag.String("stathat", "", "EZKey to use with stathat. Failure to provide will "+
		"disable stathat integration")
)

// Server provides a ServeHTTP function to enable it to be used as a
// handler. This server handles all proxying, stats collecting, etc.
type Server struct {
	syncInitOnce sync.Once

	// staticServer serves from /static/
	staticServer http.Handler

	// templates contains all the imported templates from /templates/
	templates *template.Template

	currentStatsRecords       *stats.Records
	currentStatsRecordsRWLock sync.RWMutex
}

// Init is run automatically as needed. Only the first call to it actually
// does anything. It is exposed because some clients may choose to
// explicitly initialize so that template errors are reported upon startup.
func (s *Server) Init() {
	s.syncInitOnce.Do(func() {
		s.currentStatsRecords = &stats.Records{}
		go s.statsRotator()

		// load static
		s.staticServer = http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, Prefix: ""})
		s.staticServer = http.StripPrefix("/static", s.staticServer)

		// load templates
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

func (s *Server) statsRotator() {
	if *stathatKey == "" {
		log.Println("No stathat key provided. Not posting stats to stathat.")
		return
	}
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		s.currentStatsRecordsRWLock.Lock()
		oldStatsRecords := s.currentStatsRecords
		s.currentStatsRecords = &stats.Records{}
		s.currentStatsRecordsRWLock.Unlock()
		log.Println("Rotating stats")
		for _, record := range oldStatsRecords.List() {
			postToStathat(record.Url, "Total Requests", record.TotalRequests)
		}
	}
}

func postToStathat(recordUrl, statName string, value int64) {
	log.Println("Posting into stathat!")
	stathat.PostCount(fmt.Sprint("webservice-proxy - ", statName, " - ~_all_,", recordUrl), *stathatKey, int(value))

}

func (s *Server) getCurrentStatsRecords() *stats.Records {
	s.currentStatsRecordsRWLock.RLock()
	defer s.currentStatsRecordsRWLock.RUnlock()
	return s.currentStatsRecords
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Init()
	if r.RequestURI == "/" {
		s.handleProxy(w, r)
	} else if r.RequestURI == "/admin/form/" {
		s.handleFormTest(w, r)
	} else if r.RequestURI == "/admin/" {
		s.handleAdmin(w, r)
	} else if r.RequestURI == "/stats/numClientHangups/" {
		s.handleStatsNumClientHangups(w, r)
	} else if strings.HasPrefix(r.RequestURI, "/static/") {
		s.staticServer.ServeHTTP(w, r)
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 not found")
	}
}
