package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	PongWait = 60 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

type WSClient struct {
	cons map[string]*IrcClient
	ws   *websocket.Conn
}

func (c *WSClient) HandleCommand(cmd *Command) {
	switch *cmd.Id {
	case 0:
		{

			break
		}
	case 1:
		{
			srv := *cmd.ConnectCommand.Server
			if c.cons[srv] == nil {
				cfg := &IrcConfig{
					server: srv,
					ssl:    *cmd.ConnectCommand.Ssl,
					port:   int(*cmd.ConnectCommand.Port),
				}

				nc := NewIrcClient(c, cfg)
				nc.Run()

				c.cons[srv] = nc
			}
			break
		}
	case 2:
		{
			srv := *cmd.ServerMessage.Server
			if c.cons[srv] != nil {
				c.cons[srv].SendMessage(*cmd.ServerMessage.Msg)
			}
			break
		}
	default:
		{
			stid := int32(3)
			msgtype := int32(0)
			rspm := fmt.Sprintf("Unknown command type: %v", *cmd.Id)
			rsp := &Command{
				Id: &stid,
				StatusMessage: &StatusMessage{
					Msg:     &rspm,
					Msgtype: &msgtype,
				},
			}
			c.SendMessage(rsp)
		}
	}
}

func (c *WSClient) Shutdown() {
	c.ws.Close()
	for _, v := range c.cons {
		v.Close()
	}
}

func (c *WSClient) RemoveClient(i *IrcClient) {
	srv := i.config.server
	if c.cons[srv] != nil {
		delete(c.cons, srv)
	}
}

func (c *WSClient) SendStatusMessage(tpe, code int, msg string) {
	stid := int32(3)
	st := &Command{
		Id: &stid,
		StatusMessage: &StatusMessage{
			Statuscode: &code,
			Msgtype:    &tpe,
			Msg:        &msg,
		},
	}
	c.SendMessage(st)
}

func (c *WSClient) SendMessage(cmd *Command) {
	msg, pe := proto.Marshal(cmd)
	if pe == nil {
		c.ws.WriteMessage(websocket.BinaryMessage, msg)
	}
}

func (c *WSClient) RunWS(ws *websocket.Conn) {
	c.ws = ws
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	go func() {
		defer c.Shutdown()
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
						if ume == nil {
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
		nc := &WSClient{
			cons: make(map[string]*IrcClient, 2),
		}

		nc.RunWS(c)
	} else {
		fmt.Println(err.Error())
	}
}
