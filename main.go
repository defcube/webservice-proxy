package main

import (
	"flag"
	"fmt"
	"github.com/defcube/webservice-proxy/server"
	"net/http"
)

func main() {
	port := flag.Uint64("port", 8000, "Default: 8000")
	flag.Parse()
	s := server.Server{}
	s.Init()
	http.Handle("/", &s)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%v", *port), http.DefaultServeMux)
	if err != nil {
		panic(err)
	}
}
