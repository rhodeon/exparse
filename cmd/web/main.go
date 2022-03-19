package main

import (
	"github.com/rhodeon/prettylog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/static/", serveStaticFiles)
	mux.HandleFunc("/result", calculateResult)

	server := http.Server{Addr: ":4000", Handler: mux}
	prettylog.InfoF("Starting server on %s", server.Addr)
	err := server.ListenAndServe()
	prettylog.FatalError(err)
}
