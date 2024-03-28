package model

type Token struct {
	Uid    string `json:"uid"`
	RoomId string `json:"roomId"`

	// 0: SFU; 1: Broadcaster; 2: Client
	Role int `json:"role"`

	// base64 encoded
	SDP string `json:"sdp"`
	ICE string `json:"ice"`
}

// only broadcaster request
type CreateRoomReq struct {
	Uid    string `json:"uid"`
	RoomId string `json:"roomId"`
}
