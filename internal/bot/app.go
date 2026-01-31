package bot

import (
	"encoding/json"
	"log"
	"net/url"
	"study-bot/internal/discord"
	"study-bot/internal/repository"

	"github.com/gorilla/websocket"
)

func Run() {
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

	repo := repository.Open("discord_bot:elzhqht123!@tcp(minsung.me:3306)/Discord?parseTime=true")
	sess := discord.NewSession()
	bot := NewBot(sess, repo)
	defer conn.Close()

	// 메세지 읽기
	go func() {
		defer close(sess.Stop)
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
				sess.Handshake(event)
			case 11:
				sess.NotifyAck()
			case 0:
				bot.OnEvent(event)
			}
		}
	}()

	// 메세지 쓰기
	go func() {
		for {
			select {
			case message := <-sess.Send:
				log.Printf("send: %d", message.Op)
				if err := conn.WriteJSON(message); err != nil {
					log.Printf("unmarshal error: %v", err)
					return
				}
			case <-sess.Stop:
				return
			}
		}
	}()

	<-sess.Stop
}
