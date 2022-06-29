package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	ID		string
	Conn	*websocket.Conn
	Pool	*Pool
	mu		sync.Mutex
}

func (c *Client)  Read() {
	defer func(){
		
	}()
}