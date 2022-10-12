package proxy

import (
	"net"
)

type Client struct {
	proxy      net.Conn
	connection net.Conn
}

func (client *Client) Handle() {
	go client.ReadFrom()
	go client.WriteTo()
}

func (client *Client) ReadFrom() {}

func (client *Client) WriteTo() {}
