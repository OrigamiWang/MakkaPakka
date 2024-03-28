package main

import (
	"flag"
	"log"
	"net/http"
	"signal-server/biz"
	"time"
)

var addr = flag.String("addr", ":8081", "http service address")

func main() {
	http.HandleFunc("/tk", func(w http.ResponseWriter, r *http.Request) {
		biz.TokenHandler(w, r)
	})

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
