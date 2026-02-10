package discord

import (
	"context"
	"encoding/json"
	"errors"
	"log"
)

func (c *Connection) Read(ctx context.Context) {
	defer func() {
		c.errChan <- errors.New("Conneciton Closed")
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("Read:", err)
			return
		}
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("unmarshal error: %v", err)
			c.errChan <- err
			return
		}

		log.Printf("recv: %d", event.Op)
		c.handleEvent(ctx, event)
	}
}

func (c *Connection) Send(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-c.gw.EventChannel:
			log.Printf("send: %d", message.Op)
			if err := c.ws.WriteJSON(message); err != nil {
				log.Printf("unmarshal error: %v", err)
				return
			}
		}
	}
}

func (c *Connection) handleEvent(ctx context.Context, event Event) {
	switch event.Op {
	case 0:
		if event.S != nil {
			c.gw.setSeq(*event.S)
		}
		if event.T != nil && event.D != nil && *event.T == "READY" {
			c.gw.SetResume(*event.D)
		}
		c.gw.Handler.OnEvent(event)
	case 1:
		c.SendHeartbeat()
	case 7:
		c.errChan <- errors.New("Discord Send Reconnect : opcode 7")
	case 9:
		c.HandleInvalidSession(event)
	case 10:
		c.StartHeartbeat(ctx, event)
		// 재연결
		if c.gw.isReconnect {
			c.gw.Reconnect()
		} else {
			c.gw.Handshake()
		}
	case 11:
		c.NotifyAck()
	}
}
