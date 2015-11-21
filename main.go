package main

import (
	"github.com/defcube/webservice-proxy/server"
	"net/http"
)

func main() {
	s := server.Server{}
	http.Handle("/", &s)

	http.ListenAndServe("localhost:8000", http.DefaultServeMux)
}
