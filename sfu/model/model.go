package model

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader *websocket.Upgrader

func init() {
	Upgrader = &websocket.Upgrader{
		// 允许所有CORS请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

}

type Token struct {
	Uid    string `json:"uid"`
	RoomId string `json:"roomId"`
	Role   int    `json:"role"`
	Sdp    string `json:"sdp"`
	Ice    string `json:"ice"`
}
