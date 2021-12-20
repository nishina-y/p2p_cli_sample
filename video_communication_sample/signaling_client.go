package main

import (
	"bytes"
	"log"
	"math"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

type SignalingClient struct {
	conn *websocket.Conn
}

func (c *SignalingClient) connection(addr *string, onReceive chan string) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c.conn = conn

	done := make(chan struct{})

	go func() {
		defer close(done)
		var offset int
		var receiveMessage string
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			receiveMessage += string(message)
			if bytes.Index(message, []byte("\n")) == -1 {
				continue
			}
			offset = 0
			for {
				nextOffset := strings.Index(receiveMessage[offset:], "\n")
				if nextOffset == -1 {
					break
				}
				nextOffset += offset
				onReceive <- receiveMessage[offset:nextOffset]
				offset = nextOffset + 1
			}
			receiveMessage = ""
		}
	}()
}

func (c *SignalingClient) textMessage(message string) error {
	messageByte := []byte(message)
	messageLen := len(messageByte)
	frameSize := 500
	frameCount := len(messageByte)/frameSize + 1
	for i := 0; i < frameCount; i++ {
		start := i * frameSize
		end := int(math.Min(float64(start+frameSize), float64(messageLen)))
		text := message[start:end]
		if i == frameCount-1 {
			text += "\n"
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Println("write:", err)
			return err
		}
	}
	return nil
}

func (c *SignalingClient) close() error {
	if err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Println("write close:", err)
		return err
	}
	c.conn.Close()
	return nil
}
