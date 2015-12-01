package server

import (
	"fmt"
	"github.com/bluele/gforms"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.PostForm

	// Get the url
	targetUrl := form.Get("-url")
	form.Del("-url")
	if targetUrl == "" {
		panic("missing url") // todo handle gracefully
	}

	// Look out for things closing when they shouldn't
	closeCh := w.(http.CloseNotifier).CloseNotify()
	doneCh := make(chan bool, 1)
	defer func() { doneCh <- true }()
	go func() {
		select {
		case <-closeCh:
			s.getCurrentStatsRecords().RecordClientHangup(targetUrl)
		case <-doneCh:
			// this makes it so we don't wait forever if closeCh never fires
		}
	}()

	// Get the timeout
	timeoutSecondsStr := form.Get("-timeoutSeconds")
	form.Del("-timeoutSeconds")
	if timeoutSecondsStr == "" {
		timeoutSecondsStr = "60"
	}
	timeoutSeconds, err := strconv.Atoi(timeoutSecondsStr)
	if err != nil {
		timeoutSeconds = 60
	}

	s.getCurrentStatsRecords().RecordRequest(targetUrl)

	client := http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	resp, err := client.PostForm(targetUrl, form)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			if netErr, ok := urlErr.Err.(net.Error); ok {
				if netErr.Timeout() {
					log.Println("Timeout error") // TODO count stats instead
				}
			}
		}
		w.WriteHeader(500)
		s.writeResponse(w, []byte(fmt.Sprint("Error:", err)), targetUrl)
		return
	}
	w.WriteHeader(resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) // todo
	}
	s.writeResponse(w, respBody, targetUrl)
}

func (s *Server) writeResponse(w http.ResponseWriter, respBody []byte, targetUrl string) {
	i, err := w.Write(respBody)
	if err != nil {
		log.Println("Got error writing response:", i, "Err:", err)
	}
}

func (s *Server) handleStatsNumClientHangups(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s.getCurrentStatsRecords().NumClientHangups())
}

var adminForm = gforms.DefineForm(gforms.NewFields(
	gforms.NewTextField("Foo", gforms.Validators{gforms.Required(), gforms.MaxLengthValidator(3)}),
	gforms.NewTextField("Bar", gforms.Validators{}),
))

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	recordList := s.getCurrentStatsRecords().List()
	err := s.templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"StatsRecords": recordList,
	})
	if err != nil {
		panic(err)
	}
}

func (s *Server) handleFormTest(w http.ResponseWriter, r *http.Request) {
	f := adminForm(r)
	if r.Method == "POST" {
		log.Println("Is Valid Form?", f.IsValid())
		log.Println(f.CleanedData)
	}
	err := s.templates.ExecuteTemplate(w, "formtest.html", map[string]interface{}{"Form": f})
	if err != nil {
		panic(err)
	}
}
