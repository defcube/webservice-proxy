package server_test

import (
	"fmt"
	"github.com/defcube/webservice-proxy/server"
	"github.com/defcube/webservice-proxy/server/internal/testhelpers/echohandler"
	"github.com/defcube/webservice-proxy/server/internal/testhelpers/httpclient"
	"github.com/defcube/webservice-proxy/server/internal/testhelpers/nextportgenerator"
	"github.com/defcube/webservice-proxy/server/internal/testhelpers/stoppablehttpserver"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestProxyPost(t *testing.T) {
	proxyPort, proxyUrl := nextPortAndUrl()
	targetPort, targetUrl := nextPortAndUrl()
	s := stoppablehttpserver.New(proxyPort, &server.Server{})
	defer s.Stop()
	foobarS := stoppablehttpserver.New(targetPort, &echohandler.EchoHandler{})
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

func TestTimeoutProxyPost(t *testing.T) {
	proxyPort, proxyUrl := nextPortAndUrl()
	targetPort, targetUrl := nextPortAndUrl()
	s := stoppablehttpserver.New(proxyPort, &server.Server{})
	defer s.Stop()
	targetServer := stoppablehttpserver.New(targetPort, &echohandler.EchoHandler{Delay: 5 * time.Second})
	defer targetServer.Stop()

	_, err := httpclient.HttpPostForm(proxyUrl, url.Values{
		"-url":            {targetUrl},
		"-timeoutSeconds": {"1"},
		"echoval":         {"foobar"},
	})
	assert.NotNil(t, err)
}

func TestClientHangupProxyPost(t *testing.T) {
	proxyPort, proxyUrl := nextPortAndUrl()
	targetPort, targetUrl := nextPortAndUrl()
	s := stoppablehttpserver.New(proxyPort, &server.Server{})
	defer s.Stop()
	targetServer := stoppablehttpserver.New(targetPort, &echohandler.EchoHandler{Delay: 5 * time.Second})
	defer targetServer.Stop()

	assert.Equal(t, "0", httpclient.MustHttpPostForm(proxyUrl+"/stats/numClientHangups/", url.Values{}))
	_, err := httpclient.HttpPostFormWithTimeout(proxyUrl, url.Values{
		"-url":            {targetUrl},
		"-timeoutSeconds": {"2"},
		"echoval":         {"foobar"},
	}, 1*time.Second)
	assert.NotNil(t, err)
	time.Sleep(2 * time.Second)
	assert.Equal(t, "1", httpclient.MustHttpPostForm(proxyUrl+"/stats/numClientHangups/", url.Values{}))
}

func nextPortAndUrl() (portStr string, urlStr string) {
	portStr = nextportgenerator.NextAsColonString()
	urlStr = fmt.Sprintf("http://localhost%v", portStr)
	return
}
