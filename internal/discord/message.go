package discord

import (
	"context"
	"encoding/json"
	"log"
)

func (s *Session) Read(ctx context.Context) {
	defer func() {
		s.ClosedChannel <- struct{}{}
	}()

	for {
		_, message, err := s.Conn.ReadMessage()
		if err != nil {
			log.Println("Read:", err)
			return
		}
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		log.Printf("recv: %d", event.Op)

		switch event.Op {
		case 0:
			if event.S != nil {
				s.setSeq(*event.S)
			}
			if event.T != nil && event.D != nil && *event.T == "READY" {
				s.SetResume(*event.D)
			}
			s.Handler.OnEvent(event)
		case 1:
			s.SendHeartbeat()
		case 9:
			s.Handshake()
		case 10:
			s.StartHeartbeat(ctx, event)
			// 재연결
			if s.isReconnect {
				s.Reconnect()
			} else {
				s.Handshake()
			}
			s.SendHeartbeat()
		case 11:
			s.NotifyAck()
		}
	}
}

func (s *Session) Send(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-s.EventChannel:
			log.Printf("send: %d", message.Op)
			if err := s.Conn.WriteJSON(message); err != nil {
				log.Printf("unmarshal error: %v", err)
				return
			}
		}
	}
}
