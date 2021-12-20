package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/pion/webrtc/v3"
)

const (
	DEBUG           = true
	sdpPrefix       = "@SDP:"
	candidatePrefix = "@CND:"
)

var signalingClient *SignalingClient

func debuglog(message string) {
	if DEBUG {
		fmt.Println(message)
	}
}

func main() {
	addr := flag.String("addr", "localhost:8080", "signaling server address")
	mode := flag.String("mode", "answer", "answer of offer")
	flag.Parse()
	debuglog("mode=" + *mode)

	done := make(chan bool)

	onReceive := make(chan string)

	signalingClient = &SignalingClient{}
	signalingClient.connection(addr, onReceive)

	defer signalingClient.close()

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs:           []string{"turn:50.100.100.100:3478"},
				Username:       "username",
				Credential:     "password",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := peerConnection.Close(); err != nil {
			debuglog("cannot close peerConnection: " + err.Error())
		}
	}()

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := signalCandidate(c); onICECandidateErr != nil {
			panic(onICECandidateErr)
		}
	})

	go func() {
		for {
			select {
			case message := <-onReceive:
				if strings.HasPrefix(message, sdpPrefix) {
					sdpJson := message[5:]
					sdp := webrtc.SessionDescription{}
					json.Unmarshal([]byte(sdpJson), &sdp)
					if sdpErr := peerConnection.SetRemoteDescription(sdp); sdpErr != nil {
						panic(sdpErr)
					}
					if *mode == "answer" {
						answer, err := peerConnection.CreateAnswer(nil)
						if err != nil {
							panic(err)
						}

						payload, err := json.Marshal(answer)
						if err != nil {
							panic(err)
						}
						if err := signalingClient.textMessage(sdpPrefix + string(payload)); err != nil {
							panic(err)
						}

						if err := peerConnection.SetLocalDescription(answer); err != nil {
							panic(err)
						}
					}
					candidatesMux.Lock()
					for _, c := range pendingCandidates {
						if onICECandidateErr := signalCandidate(c); onICECandidateErr != nil {
							panic(onICECandidateErr)
						}
					}
					candidatesMux.Unlock()
				} else if strings.HasPrefix(message, candidatePrefix) {
					candidate := message[5:]
					if candidateErr := peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate}); candidateErr != nil {
						panic(candidateErr)
					}
				}
			case <-done:
				return
			}
		}
	}()

	if *mode == "offer" {
		dataChannel, err := peerConnection.CreateDataChannel("data", nil)
		if err != nil {
			panic(err)
		}
		setDataChannel(dataChannel)

		offer, err := peerConnection.CreateOffer(nil)
		if err != nil {
			panic(err)
		}

		if err = peerConnection.SetLocalDescription(offer); err != nil {
			panic(err)
		}

		payload, err := json.Marshal(offer)
		if err != nil {
			panic(err)
		}
		offerJson := string(payload)
		if err = signalingClient.textMessage(sdpPrefix + offerJson); err != nil {
			panic(err)
		}
	} else if *mode == "answer" {
		peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
			debuglog(fmt.Sprintf("New DataChannel %s %d", d.Label(), d.ID()))
			setDataChannel(d)
		})
	}

	for {
	}
}

func signalCandidate(c *webrtc.ICECandidate) error {
	payload := c.ToJSON().Candidate
	if err := signalingClient.textMessage(candidatePrefix + payload); err != nil {
		return err
	}
	return nil
}

func setDataChannel(d *webrtc.DataChannel) {
	d.OnOpen(func() {
		debuglog(fmt.Sprintf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds", d.Label(), d.ID()))

		for {
			text := MustReadStdin()
			if err := d.SendText(text); err != nil {
				panic(err)
			}
		}
	})

	d.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Println("> " + string(msg.Data))
	})
}

func MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}
	return in
}
