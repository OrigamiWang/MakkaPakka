package biz

import (
	"fmt"
	"log"
	"net/http"
	"signal-server/model"
	"signal-server/util/logutil"

	"github.com/gorilla/websocket"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := model.WebSocket.Upgrager.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("handle token")
	// model.SetDefaultConn(conn)
	if err != nil {
		logutil.Info("HandleToken. read json err: %v", err.Error())
		return
	}
	clientTokenChan := make(chan *model.Token)
	errChan := make(chan int)
	go readPump(conn, clientTokenChan, errChan)
	go handleClientToken(conn, clientTokenChan, errChan)
}

func readPump(conn *websocket.Conn, tokenChan chan *model.Token, errChan chan int) {
	for {
		token := &model.Token{}
		err := conn.ReadJSON(token)
		if err != nil {
			fmt.Printf("read pump failed, re-create conneciton, err: %v\n", err)
			break
		}
		tokenChan <- token
	}
}

// 主播重连
func sendToken(conn *websocket.Conn, token *model.Token, errChan chan int) {
	err := conn.WriteJSON(token)
	if err != nil {
		fmt.Printf("send token failed, re-create conneciton, err: %v\n", err)
	}
	fmt.Println("send token...")
}

func handleSfuToken(conn *websocket.Conn, tokenChan chan *model.Token, errChan chan int) {
	for {
		token := <-tokenChan
		switch token.Role {
		case 0:
			// sfu to client or broadcaster
			fmt.Println("handle role 0")
			handleRole0(conn, token, errChan)
		}
	}
}

func handleClientToken(conn *websocket.Conn, tokenChan chan *model.Token, errChan chan int) {
	for {
		token := <-tokenChan
		switch token.Role {
		case 1:
			// broadcaster to sfu
			// create room
			fmt.Println("handle role 1")
			handleRole1(conn, token, errChan)
		case 2:
			// client to sfu
			fmt.Println("handle role 2")
			handleRole2(conn, token, errChan)
		}
	}
}

// sfu to client or broadcaster
func handleRole0(conn *websocket.Conn, token *model.Token, errChan chan int) {
	fmt.Printf("token with role 0: %+v\n", token)
	// get room
	room := model.RoomList[token.RoomId]

	// check the room's existance
	if room == nil {
		fmt.Printf("ERR: room: %s not exist\n", token.RoomId)
		return
	}

	// get client
	client := room.Clients[token.Uid]

	// check the client's existance
	if client == nil {
		fmt.Printf("ERR: client: %s not exist\n", token.Uid)
		return
	}

	// send token to client
	sendToken(client.Conn, token, errChan)
}

// broadcaster to sfu
func handleRole1(conn *websocket.Conn, token *model.Token, errChan chan int) {
	fmt.Printf("token with role 1: %+v\n", token)
	// create room
	roomId := token.RoomId
	room := model.RoomList[roomId]
	if room == nil {
		room = model.NewRoom(roomId)
		model.RoomList[roomId] = room
	}

	// check client's existance
	client := room.Clients[token.Uid]
	if client == nil {
		// add client
		client = model.NewClient(token.Uid, room, conn)
		room.AddClient(client)
	} else {
		// update client
		client.Conn = conn
	}

	// connect sfu & add listener to handler token from sfu
	if room.SfuConn == nil {
		room.SfuConn = ConnectSFU()
		sfuTokenChan := make(chan *model.Token)
		go readPump(room.SfuConn, sfuTokenChan, errChan)
		go handleSfuToken(room.SfuConn, sfuTokenChan, errChan)
	} else {
		fmt.Println("sfu already connected")
	}

	// send token to sfu
	sendToken(room.SfuConn, token, errChan)
}

// client to sfu
func handleRole2(conn *websocket.Conn, token *model.Token, errChan chan int) {
	fmt.Printf("token with role 2: %+v\n", token)
	// get room
	room := model.RoomList[token.RoomId]

	// check the room's existance
	if room == nil {
		fmt.Printf("ERR: room: %s not exist\n", token.RoomId)
		return
	}

	// get or create client
	client := room.Clients[token.Uid]
	if client == nil {
		// create client if is not exist
		client = model.NewClient(token.Uid, room, conn)
		room.Clients[token.Uid] = client
	} else {
		// update client if is exist
		client.Conn = conn
	}

	// send token to sfu
	sendToken(room.SfuConn, token, errChan)
}
