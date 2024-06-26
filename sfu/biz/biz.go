package biz

import (
	"errors"
	"fmt"
	"io"
	"sfu/model"
	"sfu/util/base64"

	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v3"
)

var PcMap = make(map[string]*webrtc.PeerConnection)
var peerConnectionConfig = webrtc.Configuration{}
var localTrackChan = make(chan *webrtc.TrackLocalStaticRTP)
var localTracks = make(map[webrtc.RTPCodecType]*webrtc.TrackLocalStaticRTP)

func HandleSFU(sdpChan, iceChan, retChan chan *model.Token) {
	// 必须先获取一个sdp
	token := <-sdpChan
	fmt.Println("get sdp token from sdpChan")
	fmt.Println(token)
	// broadcaster
	sdp := token.Sdp
	offer := webrtc.SessionDescription{}
	base64.Decode(sdp, &offer)
	fmt.Println("")

	peerConnectionConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.xten.com:3478"},
				// URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		panic(err)
	}

	// Create a InterceptorRegistry. This is the user configurable RTP/RTCP Pipeline.
	// This provides NACKs, RTCP Reports and other features. If you use `webrtc.NewPeerConnection`
	// this is enabled by default. If you are manually managing You MUST create a InterceptorRegistry
	// for each PeerConnection.
	i := &interceptor.Registry{}

	// Use the default set of Interceptors
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		panic(err)
	}

	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		panic(err)
	}
	i.Add(intervalPliFactory)

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i)).NewPeerConnection(peerConnectionConfig)
	if err != nil {
		panic(err)
	}
	PcMap[token.Uid] = peerConnection

	// Allow us to receive 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	// Allow us to receive 1 audio track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			fmt.Println("local ice candidate:")
			fmt.Println(candidate)
			ret := &model.Token{
				Uid:    token.Uid,
				RoomId: token.RoomId,
				Role:   0,
				Ice:    base64.Encode(candidate),
			}
			retChan <- ret
		}
	})
	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) { //nolint: revive
		// Create a local track, all our SFU clients will be fed via this track

		trackKind := remoteTrack.Kind().String()
		trackID := fmt.Sprintf("OrigamiWang-%s", trackKind)

		fmt.Printf("on track, kind: %s, id: %s\n", trackKind, trackID)

		localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, trackKind, trackID)
		if newTrackErr != nil {
			panic(newTrackErr)
		}
		localTrackChan <- localTrack

		rtpBuf := make([]byte, 1400)
		for {
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				panic(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err = localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				panic(err)
			}
		}
	})

	fmt.Println("Set the remote SessionDescription")
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	fmt.Println("Set the local SessionDescription")
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	go func() {
		// 不仅要收集broadcaster的ice，还要收集其他用户的ice
		for {
			token := <-iceChan
			ice := token.Ice
			var iceCandidate webrtc.ICECandidateInit
			base64.Decode(ice, &iceCandidate)
			// 需要判断一下ice是谁传过来的，然后添加到对应的pc中
			pc := PcMap[token.Uid]
			if pc == nil {
				fmt.Printf("pc is nil, uid: %v\n", token.Uid)
				continue
			} else {
				fmt.Printf("pc is exist, uid: %v\n", token.Uid)
			}
			// FIXME: ICE candidate can not set before the remote SDP set
			if err := pc.AddICECandidate(iceCandidate); err == nil {
				fmt.Println("add ice candidate")
			} else {
				fmt.Printf("add ice candidate error, err: %v\n", err)
			}
		}
	}()
	// Get the LocalDescription and take it to base64 so we can paste in browser
	// fmt.Println(base64.Encode(*peerConnection.LocalDescription()))
	//
	ret := &model.Token{
		Uid:    token.Uid,
		RoomId: token.RoomId,
		Role:   0,
		Sdp:    base64.Encode(*peerConnection.LocalDescription()),
	}
	retChan <- ret

	for {
		select {
		case localTrack := <-localTrackChan:
			localTracks[localTrack.Kind()] = localTrack
		case token := <-sdpChan:
			// 一旦有新的sdpChan过来，就代表新的client要连接webrtc，就创建一个新的peerconnnection并且存在map中，key：uid， val：pc
			fmt.Println("waiting for client pier to connect to sfu...")
			fmt.Printf("client pier token: %v", token)
			recvOnlyOffer := webrtc.SessionDescription{}
			base64.Decode(token.Sdp, &recvOnlyOffer)

			pc, err := webrtc.NewPeerConnection(peerConnectionConfig)
			if err != nil {
				panic(err)
			}
			PcMap[token.Uid] = pc
			pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
				if candidate != nil {
					fmt.Println("client pier local ice candidate:")
					fmt.Println(candidate)
					ret := &model.Token{
						Uid:    token.Uid,
						RoomId: token.RoomId,
						Role:   0,
						Ice:    base64.Encode(candidate),
					}
					retChan <- ret
				}
			})

			for _, track := range localTracks {
				rtpSender, err := pc.AddTrack(track)
				if err != nil {
					panic(err)
				}

				go func() {
					rtcpBuf := make([]byte, 1500)
					for {
						if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
							return
						}
					}
				}()
			}

			// Set the remote SessionDescription
			err = pc.SetRemoteDescription(recvOnlyOffer)
			if err != nil {
				panic(err)
			}
			fmt.Println("finish to set remote SDP")

			answer, err := pc.CreateAnswer(nil)
			if err != nil {
				panic(err)
			}

			err = pc.SetLocalDescription(answer)
			if err != nil {
				panic(err)
			}

			ret := &model.Token{
				Uid:    token.Uid,
				RoomId: token.RoomId,
				Role:   0,
				Sdp:    base64.Encode(*pc.LocalDescription()),
			}
			retChan <- ret
		}

	}
}
