/* eslint-env browser */

// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

const log = msg => {
  document.getElementById('logs').innerHTML += msg + '<br>'
}

window.createSession = isPublisher => {
  const pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.xten.com:3478'
        // urls: 'stun:stun.l.google.com:19302'
      }
    ]
  })
  pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
  pc.onicecandidate = event => {
    if (event.candidate === null) {
      document.getElementById('localSessionDescription').value = btoa(JSON.stringify(pc.localDescription))
    }
  }

  if (isPublisher) {
    navigator.mediaDevices.getUserMedia({ video: false, audio: true })
      .then(stream => {
        stream.getTracks().forEach(track => pc.addTrack(track, stream))
        document.getElementById('audio1').srcObject = stream
        pc.createOffer()
          .then(d => pc.setLocalDescription(d))
          .catch(log)
      }).catch(log)
  } else {
    pc.addTransceiver('audio')
    pc.createOffer()
      .then(d => pc.setLocalDescription(d))
      .catch(log)

    pc.ontrack = function (event) {
      const el = document.getElementById('audio1')
      el.srcObject = event.streams[0]
      el.autoplay = true
      el.controls = true
    }
  }

  window.startSession = () => {
    const sd = document.getElementById('remoteSessionDescription').value
    if (sd === '') {
      return alert('Session Description must not be empty')
    }

    try {
      pc.setRemoteDescription(JSON.parse(atob(sd)))
    } catch (e) {
      alert(e)
    }
  }

  window.copySDP = () => {
    const browserSDP = document.getElementById('localSessionDescription')

    browserSDP.focus()
    browserSDP.select()

    try {
      const successful = document.execCommand('copy')
      const msg = successful ? 'successful' : 'unsuccessful'
      log('Copying SDP was ' + msg)
    } catch (err) {
      log('Unable to copy SDP ' + err)
    }
  }

  const btns = document.getElementsByClassName('createSessionButton')
  for (let i = 0; i < btns.length; i++) {
    btns[i].style = 'display: none'
  }

  document.getElementById('signalingContainer').style = 'display: block'
}
