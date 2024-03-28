package biz

import (
	"log"

	"github.com/gorilla/websocket"
)

// 信令服务器向SFU发起websocket连接

func ConnectSFU() *websocket.Conn {
	url := "ws://127.0.0.1:8082/sfu"

	dialer := websocket.Dialer{}

	// 使用拨号器连接到服务器
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	return conn
}
