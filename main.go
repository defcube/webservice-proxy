package main
import (
	"net/http"
	"github.com/defcube/webservice-proxy/server"
)

func main() {
	s := server.Server{}
	http.ListenAndServe("localhost:8000", &s)
}

