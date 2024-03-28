package model

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var WebSocket = NewWs()

type Ws struct {
	Upgrager *websocket.Upgrader
}

func NewWs() *Ws {
	return &Ws{
		Upgrager: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
