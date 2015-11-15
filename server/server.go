package server

import (
	"io/ioutil"
	"net/http"
)

type Server struct {
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
