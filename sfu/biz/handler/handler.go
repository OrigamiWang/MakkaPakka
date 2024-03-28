package handler

import (
	"fmt"
	"net/http"
	"sfu/biz"
	"sfu/model"
)

func SfuHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("come in sfu handler...")
	conn, err := model.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading to WebSocket: %v\n", err)
		return
	}
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("ws connection close")
		return nil
	})

	sdpChan := make(chan *model.Token)
	iceChan := make(chan *model.Token)
	retChan := make(chan *model.Token)
	go func() {
		biz.HandleSFU(sdpChan, iceChan, retChan)
	}()
	go func() {
		for {
			conn.WriteJSON(<-retChan)
		}
	}()
	// 处理WebSocket连接
	for {
		token := &model.Token{}
		err := conn.ReadJSON(token)
		if err != nil {
			fmt.Printf("Error reading JSON: %v\n", err)
			break // 或者根据你的错误处理策略进行处理
		}
		if token.Sdp != "" {
			sdpChan <- token
		}
		if token.Ice != "" {
			iceChan <- token
		}
	}
}
