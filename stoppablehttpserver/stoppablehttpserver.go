package stoppablehttpserver

import (
	"net"
	"net/http"
)

type StoppableHttpServer struct {
	srv http.Server
	ln  net.Listener
}

func (s *StoppableHttpServer) Stop() {
	s.ln.Close()
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
	go ss.Start()
	return &ss
}
