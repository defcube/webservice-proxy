package server_test

import (
	"fmt"
	"github.com/defcube/webservice-proxy/server"
	"github.com/defcube/webservice-proxy/server/testhelpers/httpclient"
	"github.com/defcube/webservice-proxy/server/testhelpers/nextportgenerator"
	"github.com/defcube/webservice-proxy/server/testhelpers/stoppablehttpserver"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestProxyPost(t *testing.T) {
	proxyPort, proxyUrl := nextPortAndUrl()
	targetPort, targetUrl := nextPortAndUrl()
	s := stoppablehttpserver.New(proxyPort, &server.Server{})
	defer s.Stop()
	foobarS := stoppablehttpserver.New(targetPort, &EchoHandler{})
	defer foobarS.Stop()

	r, err := httpclient.HttpPostForm(proxyUrl, url.Values{
		"-url":    {targetUrl},
		"echoval": {"foobar"},
	})
	if err != nil {
		panic(err)
	}
	if r != "foobar" {
		t.Fatal("Unexpected response:", r)
	}
}

func nextPortAndUrl() (portStr string, urlStr string) {
	portStr = nextportgenerator.NextAsColonString()
	urlStr = fmt.Sprintf("http://localhost%v", portStr)
	return
}

func TestTimeoutProxyPost(t *testing.T) {
	proxyPort, proxyUrl := nextPortAndUrl()
	targetPort, targetUrl := nextPortAndUrl()
	s := stoppablehttpserver.New(proxyPort, &server.Server{})
	defer s.Stop()
	targetServer := stoppablehttpserver.New(targetPort, &EchoHandler{Delay: 5 * time.Second})
	defer targetServer.Stop()

	r, err := httpclient.HttpPostForm(proxyUrl, url.Values{
		"-url":    {targetUrl},
		"echoval": {"foobar"},
	})
	if err != nil {
		panic(err)
	}
	// TODO expect an error here instead
	assert.Equal(t, r, "foobar")
}

type EchoHandler struct {

	// Delay specifies how much to sleep before serving the response
	Delay time.Duration
}

func (s *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(s.Delay)
	fmt.Fprint(w, r.PostFormValue("echoval"))
}
