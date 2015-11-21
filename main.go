package main

import (
	"github.com/defcube/webservice-proxy/server"
	"github.com/defcube/webservice-proxy/server/static"
	"github.com/elazarl/go-bindata-assetfs"
	"net/http"
)

func main() {
	s := server.Server{}
	http.Handle("/", &s)

	fs := http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, Prefix: ""})
	fs = http.StripPrefix("/static", fs)
	http.Handle("/static/", fs)

	http.ListenAndServe("localhost:8000", http.DefaultServeMux)
}
