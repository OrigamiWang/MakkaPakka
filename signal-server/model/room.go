package model

import (
	"github.com/gorilla/websocket"
)

var RoomList = make(map[string]*Room)

type Room struct {
	Id string // 房间号
	// broadcaster *Broadcaster       // 主播
	// sfu     *SFU               // Selected Forwarding Unit
	Clients map[string]*Client // 观众
	SfuConn *websocket.Conn    // 信令服作为ws客户端向SFU发起的连接
}

func NewRoom(id string) *Room {
	return &Room{
		Id:      id,
		Clients: make(map[string]*Client),
	}
}

func (room *Room) AddClient(client *Client) {
	room.Clients[client.Id] = client
}
