/* eslint-env browser */

// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

const log = msg => {
    document.getElementById('logs').innerHTML += msg + '<br>'
}

var isPublisher

window.createSession = isPublisherTemp => {
    isPublisher = isPublisherTemp
    const btns = document.getElementsByClassName('createSessionButton')
    for (let i = 0; i < btns.length; i++) {
        btns[i].style = 'display: none'
    }

    document.getElementById('signalingContainer').style = 'display: block'
}



window.startSession = () => {
    // 先输入uid, roomId, role
    var uid = document.getElementById('uid').value
    var roomId = document.getElementById('roomId').value
    var role = parseInt(document.getElementById('role').value)
    if (uid == '' || roomId == '' || role == '') {
        console.log('uid, roomId, role cannot be empty')
        return
    }

    console.log("uid: " + uid + ", roomId: " + roomId + ", role: " + role)
    console.log("isPublisher: " + isPublisher);
    var ws = new WebSocket("ws://127.0.0.1:8081/tk");
    // var ws = new WebSocket("ws://192.168.1.126:8081/tk");
    console.log("ws state: " + ws.readyState);
    ws.onclose = function (evt) {
        console.log("ws close");
    };
    ws.onmessage = function (evt) {
        var token = evt.data;
        var sdp = JSON.parse(token)["sdp"]
        var ice = JSON.parse(token)["ice"]

        // add remote sdp
        if (sdp != "") {
            document.getElementById('remoteSessionDescription').value = sdp
            try {
                console.log('get remote SDP');
                pc.setRemoteDescription(JSON.parse(atob(sdp)))
            } catch (e) {
                alert(e)
            }
        }

        // add remote ice
        if (ice != "") {
            pc.addIceCandidate(JSON.parse(atob(ice))).then(() => {
                console.log('get remote ice candidate:');
                console.log(JSON.parse(atob(ice)));
            }).catch(error => {
                console.error('Error adding received ICE candidate', error);
            });
        }

    };
    const pc = new RTCPeerConnection({
        iceServers: [
            {
                urls: 'stun:stun.xten.com:3478'
                // urls: 'stun:stun.l.google.com:19302'
            }
        ]
    })
    ws.onopen = function (evt) {
        pc.oniceconnectionstatechange = e => {
            log(pc.iceConnectionState)
        }
        pc.onicecandidate = event => {
            if (event.candidate != null) {
                ice = btoa(JSON.stringify(event.candidate))
                var token = {
                    uid: uid,
                    roomId: roomId,
                    role: role,
                    sdp: "",
                    ice: ice,
                };
                console.log("get local ice candidate: ");
                console.log(ice);
                // 将对象转换为JSON字符串
                var jsonToken = JSON.stringify(token);
                if (ws.readyState == WebSocket.OPEN) {
                    ws.send(jsonToken)
                } else {
                    console.log("ws is not open");
                }
            } else {
                console.log("candidate is gathered");
            }
        }

        if (isPublisher) {
            navigator.mediaDevices.getUserMedia({ video: true, audio: true })
                .then(stream => {
                    stream.getTracks().forEach(track => {
                        pc.addTrack(track, stream)
                        console.log("track kind: " + track.kind);
                        if (track.kind == 'video') {
                            const videoElement = document.getElementById('video1');
                            videoElement.srcObject = stream;
                            videoElement.autoplay = true;
                            videoElement.controls = true;
                            // videoElement.muted = false
                        } else if (track.kind == 'audio') {
                            const audioElement = document.getElementById('audio1')
                            audioElement.srcObject = stream;
                            audioElement.autoplay = true;
                            audioElement.controls = true;
                            audioElement.muted = true
                        }
                    })
                    // const videoElement = document.getElementById('video1');
                    // videoElement.srcObject = stream;
                    // videoElement.autoplay = true;
                    // videoElement.controls = true;
                    // videoElement.muted = true

                    // here may badly create offer twice
                    console.log("create offer.......");
                    pc.createOffer()
                        .then(sdp => {
                            console.log("get local SDP:");
                            console.log(sdp);
                            pc.setLocalDescription(sdp)
                            document.getElementById('localSessionDescription').value = btoa(JSON.stringify(sdp))
                            var localSdp = btoa(JSON.stringify(sdp))
                            // 可以在这里把sdp传给信令服
                            // connect ws
                            var token = {
                                uid: uid,
                                roomId: roomId,
                                role: role,
                                sdp: localSdp,
                                ice: "",
                            };
                            console.log(token);
                            // 将对象转换为JSON字符串
                            var jsonToken = JSON.stringify(token);
                            if (ws.readyState == WebSocket.OPEN) {
                                ws.send(jsonToken)
                            } else {
                                console.log("ws is not open");
                            }
                        })
                        .catch(log)
                }).catch(log)
        } else {
            pc.addTransceiver('audio')
            pc.addTransceiver('video')
            // create offer
            pc.createOffer()
                .then(sdp => {
                    console.log("get local SDP:");
                    console.log(sdp);
                    pc.setLocalDescription(sdp)
                    document.getElementById('localSessionDescription').value = btoa(JSON.stringify(sdp))
                    var localSdp = btoa(JSON.stringify(sdp))
                    // 可以在这里把sdp传给信令服
                    // connect ws
                    var token = {
                        uid: uid,
                        roomId: roomId,
                        role: role,
                        sdp: localSdp,
                        ice: "",
                    };
                    console.log(token);
                    // 将对象转换为JSON字符串
                    var jsonToken = JSON.stringify(token);
                    if (ws.readyState == WebSocket.OPEN) {
                        ws.send(jsonToken)
                    } else {
                        console.log("ws is not open");
                    }
                })

            // get track
            pc.ontrack = function (event) {
                event.streams.forEach(stream => {
                    stream.getTracks().forEach(track => {
                        console.log("track kind: " + track.kind);
                        if (track.kind == 'video') {
                            const videoElement = document.getElementById('video1');
                            videoElement.srcObject = stream;
                            videoElement.autoplay = true;
                            videoElement.controls = true;
                            // videoElement.muted = false
                        } else if (track.kind == 'audio') {
                            const audioElement = document.getElementById('audio1')
                            audioElement.srcObject = stream;
                            audioElement.autoplay = true;
                            audioElement.controls = true;
                            audioElement.muted = false
                        }
                    })
                    console.log("stream id: " + stream.id);
                    // 设置视频元素的srcObject来播放视频和音频
                    
                })
            }
        }
    }

}