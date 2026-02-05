package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	Interval       int
	Send           chan Event
	lastAckReceive bool
	seq            *int
	mu             sync.Mutex
}

func (s *Session) setSeq(seq int) {
	defer s.mu.Unlock()
	s.mu.Lock()

	fmt.Printf(">>> %v", seq)
	s.seq = &seq
}

func (s *Session) getSeq() *int {
	defer s.mu.Unlock()
	s.mu.Lock()

	return s.seq
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
		lastAckReceive: false,
	}

	return &session
}

func (s *Session) Open(ctx context.Context, handler func(Event)) {
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Open Discord Gateway...")
	u := url.URL{
		Scheme:   "wss",
		Host:     "gateway.discord.gg",
		RawQuery: "v=10&encoding=json",
	}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial : ", err)
	}
	defer conn.Close()

	// 메세지 읽기
	go func() {
		for {
			_, message, err := conn.ReadMessage()
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
			case 1:
				s.SendHeartbeat()
			case 10:
				s.Handshake(innerCtx, event)
			case 11:
				s.NotifyAck()
			case 0:
				if event.S != nil {
					s.setSeq(*event.S)
				}
				handler(event)
			}
		}
	}()

	// 메세지 쓰기
	for {
		select {
		case <-innerCtx.Done():
			return
		case message := <-s.Send:
			log.Printf("send: %d", message.Op)
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("unmarshal error: %v", err)
				return
			}
		}
	}
}

func (s *Session) Handshake(ctx context.Context, event Event) {
	var handshakeEvent HandshakePayload
	if err := json.Unmarshal(*event.D, &handshakeEvent); err != nil {
		log.Printf("Handshake marshal error %v", err)
		return
	}

	log.Printf("Start heartbeat interval %d\n", handshakeEvent.HeartbeatInterval)
	s.Interval = handshakeEvent.HeartbeatInterval
	s.StartHeartbeat(ctx)

	identifyPayload, err := json.Marshal(IdentifyPayload{
		Token: os.Getenv("DISCORD_BOT_TOKEN"),
		Properties: struct {
			Os      string "json:\"os\""
			Browser string "json:\"browser\""
			Device  string "json:\"device\""
		}{
			Os:      runtime.GOOS,
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
		S:  s.seq,
	}

	s.SendHeartbeat()
}

func (s *Session) SendHeartbeat() {
	s.Send <- Event{
		Op: 1,
		S:  s.seq,
	}
}

func (s *Session) StartHeartbeat(ctx context.Context) {
	t := time.NewTicker(time.Duration(s.Interval) * time.Millisecond)
	go func() {
		for {
			select {
			case <-t.C:
				if !s.lastAckReceive {
					log.Println("ack is dead.")
					return
				}
				s.lastAckReceive = false
				log.Printf("heartbeat duration: %d", time.Duration(s.Interval))
				s.Send <- Event{Op: 1}
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}()
}

func (s *Session) NotifyAck() {
	s.setAck(true)
}
