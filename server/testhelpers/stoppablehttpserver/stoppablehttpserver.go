package stoppablehttpserver

import (
	"net"
	"net/http"
)

type StoppableHttpServer struct {
	srv  http.Server
	ln   net.Listener
	done chan bool
}

func (s *StoppableHttpServer) Stop() {
	s.ln.Close()
	<-s.done //wait until server shuts down
}

func (s *StoppableHttpServer) Start() error {
	return s.srv.Serve(s.ln)
}

// NewStoppableHttpServer creates a StoppableHttpServer and calls Start() on it
func New(laddr string, handler http.Handler) *StoppableHttpServer {
	ln, err := net.Listen("tcp", laddr)
	if err != nil {
		panic(err)
	}
	s := http.Server{Handler: handler}
	ss := StoppableHttpServer{srv: s, ln: ln}

	// Start the goroutine, and then AFTER make the ss.done channel to try to avoid a possible
	// race condition where the server isn't started and the test expects it to be.
	go func() {
		ss.Start()
		ss.done <- true
	}()
	ss.done = make(chan bool, 1)

	return &ss
}
