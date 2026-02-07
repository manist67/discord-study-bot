package discord

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Handler interface {
	OnEvent(Event)
}
type Session struct {
	Conn    *websocket.Conn
	Handler Handler

	Interval       int
	EventChannel   chan Event
	ClosedChannel  chan struct{}
	lastAckReceive bool
	mu             sync.Mutex

	isReconnect   bool
	seq           int
	connectionURL string
	sessionId     string
}

func (s *Session) setResumeValue(resultURL string, sessionId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	decodeString, err := url.QueryUnescape(resultURL)
	if err != nil {
		log.Fatalf("Fail to decode string %s", decodeString)
		return
	}

	s.isReconnect = true
	s.connectionURL = strings.Replace(decodeString, "wss://", "", 0)
	s.sessionId = sessionId
}

func (s *Session) setSeq(seq int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq = seq
}

func (s *Session) getSeq() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.seq
}

func (s *Session) setAck(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastAckReceive = v
}

func (s *Session) getAck() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.lastAckReceive
}

func NewSession() *Session {
	session := Session{
		Interval:       -1,
		EventChannel:   make(chan Event, 10),
		ClosedChannel:  make(chan struct{}),
		lastAckReceive: false,

		// isReconnect:   true,
		// connectionURL: "gateway-us-east1-d.discord.gg",
		// seq:           5,
		// sessionId:     "278ad825418529133053e38fbf2d78db",
	}

	return &session
}

func (s *Session) Open(ctx context.Context, handler Handler) {
	s.Handler = handler
	innerCtx, cancel := context.WithCancel(ctx)

	for {
		var u url.URL
		if !s.isReconnect {
			log.Printf("Open Discord Gateway...")
			u = url.URL{
				Scheme:   "wss",
				Host:     "gateway.discord.gg",
				RawQuery: "v=10&encoding=json",
			}
		} else {
			log.Printf("Reconnect Discord Gateway...")
			u = url.URL{
				Scheme:   "wss",
				Host:     s.connectionURL,
				RawQuery: "v=10&encoding=json",
			}
		}
		log.Printf("connecting to %s", u.String())

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		s.Conn = conn
		if err != nil {
			log.Fatal("dial : ", err)
		}

		// 메세지 읽기
		go s.Read(innerCtx)

		// 메세지 쓰기
		go s.Send(innerCtx)

		select {
		case <-s.ClosedChannel:
			s.Conn.Close()
			continue
		case <-innerCtx.Done():
			log.Printf("%s %s %d", s.connectionURL, s.sessionId, s.seq)
			s.Conn.Close()
			cancel()
			return
		}
	}
}

func (s *Session) Handshake() {
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
	s.EventChannel <- Event{
		Op: 2,
		D:  &raw,
		S:  &s.seq,
	}
}

func (s *Session) Reconnect() {
	resumePayload, err := json.Marshal(ResumePayload{
		Token:     os.Getenv("DISCORD_BOT_TOKEN"),
		SessionId: s.sessionId,
		Seq:       s.seq,
	})
	if err != nil {
		log.Fatal("Fail to Reconnect : ", err)
	}
	d := json.RawMessage(resumePayload)
	s.EventChannel <- Event{
		Op: 6,
		S:  &s.seq,
		D:  &d,
	}
}

func (s *Session) SetResume(p json.RawMessage) {
	var payload ReadyPayload
	if err := json.Unmarshal(p, &payload); err != nil {
		log.Printf("Fail to unmarshal payload %s", string(p))
		return
	}

	s.setResumeValue(payload.ResumeGatewayUrl, payload.SessionId)
}

func (s *Session) SendHeartbeat() {
	s.EventChannel <- Event{
		Op: 1,
		S:  &s.seq,
	}
}

func (s *Session) StartHeartbeat(ctx context.Context, event Event) {
	var handshakeEvent HandshakePayload
	if err := json.Unmarshal(*event.D, &handshakeEvent); err != nil {
		log.Printf("Handshake marshal error %v", err)
		return
	}

	log.Printf("Start heartbeat interval %d\n", handshakeEvent.HeartbeatInterval)
	s.Interval = handshakeEvent.HeartbeatInterval

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
				s.SendHeartbeat()
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
