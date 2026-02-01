package bot

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"study-bot/internal/discord"
	"study-bot/internal/repository"

	"github.com/gorilla/websocket"
)

type Bot struct {
	session       *discord.Session
	repo          *repository.Conn
	applicationId string
}

func NewBot(r *repository.Conn) *Bot {
	sess := discord.NewSession()

	return &Bot{session: sess, repo: r}
}

func (b *Bot) Run(ctx context.Context) {

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

			var event discord.Event
			if err := json.Unmarshal(message, &event); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			if event.T != nil {
				log.Printf("recv: %d %s", event.Op, *event.T)
			}

			switch event.Op {
			case 10:
				b.session.Handshake(ctx, event)
			case 11:
				b.session.NotifyAck()
			case 0:
				b.OnEvent(event)
			}
		}
	}()

	// 메세지 쓰기
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-b.session.Send:
				log.Printf("send: %d", message.Op)
				if err := conn.WriteJSON(message); err != nil {
					log.Printf("unmarshal error: %v", err)
					return
				}
			}

		}
	}()

	<-ctx.Done()
	log.Println("Terminating Discord bot")
	conn.Close()
	log.Println("Discord Bot terminated...")
}
