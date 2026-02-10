package discord

import (
	"context"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Handler interface {
	OnEvent(Event)
}
type Gateway struct {
	Handler Handler
	mu      sync.Mutex

	EventChannel chan Event

	isReconnect   bool
	seq           int
	connectionURL string
	sessionId     string
}

func (s *Gateway) setResumeValue(resultURL string, sessionId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	decodeString, err := url.PathUnescape(resultURL)
	if err != nil {
		log.Fatalf("Fail to decode string %s", decodeString)
		return
	}

	s.isReconnect = true
	s.connectionURL = strings.Replace(decodeString, "wss://", "", 1)
	s.sessionId = sessionId
}

func (s *Gateway) setSeq(seq int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq = seq
}

func (s *Gateway) getSeq() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.seq
}

func NewSession() *Gateway {
	session := Gateway{
		EventChannel: make(chan Event, 10),
	}

	return &session
}

func (s *Gateway) getUrl() url.URL {
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
	return u
}

func (s *Gateway) Open(ctx context.Context, handler Handler) {
	s.Handler = handler

	for {
		u := s.getUrl()
		log.Printf("connecting to %s", u.String())
		innerCtx, cancel := context.WithCancel(ctx)

		ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		connection := Connection{
			gw:      s,
			ws:      ws,
			errChan: make(chan error),
		}
		if err != nil {
			log.Fatal("dial : ", err)
		}

		// 메세지 읽기
		go connection.Read(innerCtx)

		// 메세지 쓰기
		go connection.Send(innerCtx)

		select {
		case <-connection.errChan:
			cancel()
			connection.ws.Close()
			continue

		case <-ctx.Done():
			log.Printf("%s %s %d", s.connectionURL, s.sessionId, s.seq)
			cancel()
			connection.ws.Close()
			return
		}
	}
}
