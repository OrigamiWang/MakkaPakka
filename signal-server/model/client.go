package model

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   string
	Room *Room
	Conn *websocket.Conn
}

// type Broadcaster struct {
// 	room *Room
// 	conn *websocket.Conn
// }

type SFU struct {
	room *Room
	conn *websocket.Conn
}

func NewClient(id string, room *Room, conn *websocket.Conn) *Client {
	return &Client{
		Id:   id,
		Room: room,
		Conn: conn,
	}
}

// func NewBroadcaster(room *Room, conn *websocket.Conn) *Broadcaster {
// 	return &Broadcaster{
// 		room: room,
// 		conn: conn,
// 	}
// }
