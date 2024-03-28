package biz

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"signal-server/model"
// 	"signal-server/util/logutil"

// 	"github.com/gorilla/websocket"
// )

// func HandleToken(room *model.Room, w http.ResponseWriter, r *http.Request) {
// 	conn, err := model.WebSocket.Upgrager.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	fmt.Println("handle token")
// 	// model.SetDefaultConn(conn)
// 	if err != nil {
// 		logutil.Info("HandleToken. read json err: %v", err.Error())
// 		return
// 	}
// 	go readPump(room, conn)
// }

// // FIXME: on websocket close
// func readPump(room *model.Room, conn *websocket.Conn) {
// 	for {
// 		token := &model.Token{}
// 		err := conn.ReadJSON(token)
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				fmt.Printf("ws unexpected close, error: %v\n", err)
// 			} else {
// 				fmt.Printf("ws closed error: %v\n", err)
// 			}
// 			// 案例来说，ws退出，需要关闭client实例内的ws连接
// 			break // 退出循环，因为连接已经关闭或发生错误
// 		}
// 		if token.Role == 1 || token.Role == 2 {
// 			// broadcaster or client to sfu
// 			if room.Clients[token.Uid] == nil {
// 				// create client if is nil
// 				fmt.Println("the client is nil, create client")
// 				room.Clients[token.Uid] = model.NewClient(room, conn)
// 			} else {
// 				// check the client ws connection is closed
// 				pingToken := &model.Token{}
// 				err := room.Clients[token.Uid].Conn.WriteJSON(pingToken)
// 				if err != nil {
// 					fmt.Println("re-create ws connection")
// 					room.Clients[token.Uid].Conn = conn
// 				}
// 			}
// 			fmt.Println("send to sfu")
// 			err := room.SfuConn.WriteJSON(token)
// 			if err != nil {
// 				fmt.Printf("sfu ws connection closed, re-create conneciton")
// 				room.SfuConn = ConnectSFU()
// 			}
// 		} else if token.Role == 0 {
// 			// sfu to client
// 			if client_conn := room.Clients[token.Uid]; client_conn != nil {
// 				// 判断连接是否关闭
// 				err := client_conn.Conn.WriteJSON(token)
// 				if err != nil {
// 					fmt.Printf("send to client error, err: %v\n", err)
// 				}
// 			}
// 		}
// 	}
// }

// func ListenSfu(room *model.Room) {
// 	room.SfuConn = ConnectSFU()
// 	go readPump(room, room.SfuConn)
// }
