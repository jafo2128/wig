package main;

import (
	"net"
	"fmt"
	"bufio"
	"crypto/tls"
)

type IrcConfig struct {
	server string
	port int
	ssl bool
}

type IrcClient struct {
	config *IrcConfig
	conn *net.Conn
	conn_tls *tls.Conn
	
	io *bufio.ReadWriter
}

func (i *IrcClient) Run(cli *Client){
	fmt.Println("Starting new irc connection to:", i.config.server)
	if i.config.ssl {
		tc := &tls.Config { InsecureSkipVerify: true, }
		nc, er := tls.Dial("tcp", fmt.Sprintf("%v:%v", i.config.server, i.config.port), tc)
		if er == nil {
			i.conn_tls = nc
			
			stid := int32(3)
			fwm := &Command{
				Id: &cmdid,
				StatusMessage: &StatusMessage {
					Server: &i.config.server,
					Msg: &s,
				},
			}
			cli.SendMessage(fwm)
						
			i.io = bufio.NewReadWriter(bufio.NewReader(i.conn_tls), bufio.NewWriter(i.conn_tls))
			go func() {
				defer i.conn_tls.Close()
				cmdid := int32(2)
				for {
					s, rer := i.io.ReadString('\n')
					if rer == nil {
						fwm := &Command{
							Id: &cmdid,
							ServerMessage: &ServerMessage {
								Server: &i.config.server,
								Msg: &s,
							},
						}
						cli.SendMessage(fwm)
					}else{
						fmt.Println("Read error:", rer.Error())
						break
					}
				}
			}()
		}else{
			fmt.Println("Connect error:", er.Error())
		}
	}
}

func NewIrcClient(cfg *IrcConfig) *IrcClient {
	return &IrcClient{
		config: cfg,
	}
}