package server_test

import (
	"fmt"
	"github.com/defcube/webservice-proxy/httpclient"
	"github.com/defcube/webservice-proxy/server"
	"github.com/defcube/webservice-proxy/stoppablehttpserver"
	"net/http"
	"net/url"
	"testing"
)

func TestProxyPost(t *testing.T) {
	s := stoppablehttpserver.New(":8098", &server.Server{})
	defer s.Stop()
	foobarS := stoppablehttpserver.New(":8099", &EchoHandler{})
	defer foobarS.Stop()

	r, err := httpclient.HttpPostForm("http://localhost:8098", url.Values{
		"-url":    {"http://localhost:8099"},
		"echoval": {"foobar"},
	})
	if err != nil {
		panic(err)
	}
	if r != "foobar" {
		t.Fatal("Unexpected response:", r)
	}
}

type EchoHandler struct {
}

func (s *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.PostFormValue("echoval"))
}
