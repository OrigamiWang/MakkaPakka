package main

import (
	"log"
	"net/http"
	"sfu/biz/handler"
)

func main() {
	// 设置WebSocket路由
	http.HandleFunc("/sfu", handler.SfuHandler)

	// 启动HTTP服务器
	log.Println("Starting server on :8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
