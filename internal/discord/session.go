package discord

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Session struct {
	Interval       int
	Send           chan Event
	Stop           chan struct{}
	lastAckReceive bool
	mu             sync.Mutex
}

func (s *Session) setAck(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastAckReceive = v
}

func (s *Session) getAck() bool {
	defer s.mu.Unlock()
	s.mu.Lock()
	return s.lastAckReceive
}

func NewSession() *Session {
	session := Session{
		Interval:       -1,
		Send:           make(chan Event, 10),
		Stop:           make(chan struct{}),
		lastAckReceive: false,
	}

	return &session
}

func (s *Session) Handshake(event Event) {
	var handshakeEvent HandshakePayload
	if err := json.Unmarshal(*event.D, &handshakeEvent); err != nil {
		log.Printf("Handshake marshal error %v", err)
		return
	}

	log.Printf("Start heartbeat interval %d\n", handshakeEvent.HeartbeatInterval)
	s.Interval = handshakeEvent.HeartbeatInterval
	s.StartHeartbeat()

	identifyPayload, err := json.Marshal(IdentifyPayload{
		Token: os.Getenv("DISCORD_BOT_TOKEN"),
		Properties: struct {
			Os      string "json:\"os\""
			Browser string "json:\"browser\""
			Device  string "json:\"device\""
		}{
			Os:      "windows",
			Browser: "discordbot",
			Device:  "discordbot",
		},
		Intents: 1<<7 | 1<<0,
	})
	if err != nil {
		log.Printf("Handshake marshal error %v", err)
		return
	}

	raw := json.RawMessage(identifyPayload)
	s.Send <- Event{
		Op: 2,
		D:  &raw,
	}

	s.Send <- Event{Op: 1}
}

func (s *Session) StartHeartbeat() {
	t := time.NewTicker(time.Duration(s.Interval) * time.Millisecond)
	go func() {
		for {
			select {
			case <-t.C:
				if !s.lastAckReceive {
					log.Println("ack is dead.")
					s.Stop <- struct{}{}
					panic(true)
				}
				s.lastAckReceive = false
				log.Printf("heartbeat duration: %d", time.Duration(s.Interval))
				s.Send <- Event{Op: 1}
			case <-s.Stop:
				return
			}
		}
	}()
}

func (s *Session) NotifyAck() {
	s.setAck(true)
}
