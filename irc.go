package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/golang/protobuf/proto"
)

const (
	PongWait = 60 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
	
}

type Client struct {
	cons map[string]*IrcClient
	ws   *websocket.Conn
}

func (c *Client) HandleCommand(cmd *Command) {
	switch(*cmd.Id) {
		case 0: {
			
			break
		}
		case 1: {
			srv := *cmd.ConnectCommand.Server
			if c.cons[srv] == nil {
				cfg := &IrcConfig {
					server: srv,
					ssl: *cmd.ConnectCommand.Ssl,
					port: int(*cmd.ConnectCommand.Port),
				}
				
				nc := NewIrcClient(cfg)
				nc.Run(c)
				
				c.cons[srv] = nc
			}
			break 
		}
		default: {
			msgtype := int32(0)
			rspm := fmt.Sprintf("Unknown command type: %v", cmd.Id)
			fmt.Println(rspm, cmd.Id)
			rsp := &Command{
				Id: cmd.Id,
				StatusMessage: &StatusMessage {
					Msg: &rspm,
					Msgtype: &msgtype,
				},
			}
			c.SendMessage(rsp)
		}
	}
}

func (c *Client) SendMessage(cmd *Command){
	msg, pe := proto.Marshal(cmd)
	if pe == nil {
		c.ws.WriteMessage(websocket.BinaryMessage, msg)
	}
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
				case websocket.BinaryMessage:
					{
						cmd := new(Command)
						ume := proto.Unmarshal(message, cmd)
						if ume == nil{
							fmt.Println("Got msg: ", cmd)
							c.HandleCommand(cmd)
						}
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
	w.Header().Set("Sec-Websocket-Protocol", "irc")
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err == nil {
		nc := &Client{
			cons: make(map[string]*IrcClient, 2),
		}

		nc.RunWS(c)
	} else {
		fmt.Println(err.Error())
	}
}
