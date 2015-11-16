package main

import (
	"github.com/defcube/webservice-proxy/server"
	"net/http"
)

func main() {
	s := server.Server{}
	http.ListenAndServe("localhost:8000", &s)
}
