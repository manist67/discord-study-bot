package discord

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	gw *Gateway
	ws *websocket.Conn
	mu sync.Mutex

	errChan chan error

	lastAckReceive bool
	Interval       int
}

func (c *Connection) setAck(v bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastAckReceive = v
}

func (c *Connection) getAck() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.lastAckReceive
}

func (c *Connection) NotifyAck() {
	c.setAck(true)
}
