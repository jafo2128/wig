package main

import (
	"fmt"
	"net/http"
	"time"

	irc "github.com/fluffle/goirc/client"
	"github.com/gorilla/websocket"
)

/*"github.com/golang/protobuf/proto"*/

const (
	PongWait = 60 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

type Client struct {
	conn *irc.Conn
	ws   *websocket.Conn
}

func (c *Client) RunWS(ws *websocket.Conn) {
	c.ws = ws
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	go func() {
		defer c.ws.Close()
		for {
			mt, message, err := c.ws.ReadMessage()
			if err == nil {
				switch mt {
				case websocket.TextMessage:
					{
						fmt.Println(string(message))
						break
					}
				}
			} else {
				fmt.Println(err.Error())
				break
			}
		}
	}()
}

func NewClient(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		nc := &Client{}
		nc.RunWS(c)
	} else {
		fmt.Println(err.Error())
	}
}
