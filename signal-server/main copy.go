package main

// import (
// 	"flag"
// 	"log"
// 	"net/http"
// 	"signal-server/biz"
// 	"signal-server/model"
// 	"time"
// )

// var addr = flag.String("addr", ":8081", "http service address")

// func main() {
// 	room := model.NewRoom()
// 	biz.ListenSfu(room)
// 	http.HandleFunc("/tk", func(w http.ResponseWriter, r *http.Request) {
// 		biz.HandleToken(room, w, r)
// 	})

// 	server := &http.Server{
// 		Addr:              *addr,
// 		ReadHeaderTimeout: 3 * time.Second,
// 	}
// 	log.Fatal(server.ListenAndServe())
// }
