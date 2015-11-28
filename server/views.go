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
		fmt.Fprint(w, "Error:", err)
		return
	}
	w.WriteHeader(resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) // todo
	}
	s.statsRecords.RecordRequest(targetUrl)
	w.Write(respBody)
}

var adminForm = gforms.DefineForm(gforms.NewFields(
	gforms.NewTextField("Foo", gforms.Validators{gforms.Required(), gforms.MaxLengthValidator(3)}),
	gforms.NewTextField("Bar", gforms.Validators{}),
))

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	recordList := s.statsRecords.List()
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
