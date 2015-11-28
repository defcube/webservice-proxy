package echohandler

import (
	"fmt"
	"net/http"
	"time"
)

type EchoHandler struct {

	// Delay specifies how much to sleep before serving the response
	Delay time.Duration
}

func (s *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(s.Delay)
	fmt.Fprint(w, r.PostFormValue("echoval"))
}
