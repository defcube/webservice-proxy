package server

import (
	"github.com/bluele/gforms"
	"io/ioutil"
	"log"
	"net/http"
)

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
	s.statsRecords.RecordRequest(url)
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
