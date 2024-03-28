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
    var ws = new WebSocket("ws://9.134.75.180:8081/tk");
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
                urls: 'stun:stun.l.google.com:19302'
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
            navigator.mediaDevices.getUserMedia({ video: false, audio: true })
                .then(stream => {
                    stream.getTracks().forEach(track => pc.addTrack(track, stream))
                    document.getElementById('audio1').srcObject = stream
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
                const el = document.getElementById('audio1')
                el.srcObject = event.streams[0]
                el.autoplay = true
                el.controls = true
            }
        }
    }

}