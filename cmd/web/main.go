package main

import (
	"flag"
	"github.com/rhodeon/prettylog"
	"net/http"
	"os"
)

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "4000"
	}
	addr := flag.String("addr", ":"+port, "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/static/", serveStaticFiles)
	mux.HandleFunc("/result", calculateResult)

	server := http.Server{Addr: *addr, Handler: mux}
	prettylog.InfoF("Starting server on %s", server.Addr)
	err := server.ListenAndServe()
	prettylog.FatalError(err)
}
