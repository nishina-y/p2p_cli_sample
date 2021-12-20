package main

import (
	"bytes"
	"log"
	"net/url"

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
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			offset = 0
			for {
				nextOffset := bytes.Index(message[offset:], []byte("\n"))
				if nextOffset == -1 {
					break
				}
				nextOffset += offset
				messageStr := string(message[offset:nextOffset])

				onReceive <- messageStr

				offset = nextOffset + 1
			}
		}
	}()
}

func (c *SignalingClient) textMessage(message string) error {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(message+"\n")); err != nil {
		log.Println("write:", err)
		return err
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
