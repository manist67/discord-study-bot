package bot

import (
	"context"
	"log"
	"study-bot/internal/discord"
	"study-bot/internal/repository"
	"time"
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
	restartTime := 2 * time.Second
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping bot...")
			return
		default:
			b.session.Open(ctx, b.OnEvent)

			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping bot...")
				return
			case <-time.After(restartTime):
				log.Println("Restart discord bot")
				restartTime *= 2
				if restartTime > 30*time.Second {
					log.Printf("Max retry interval reched. Shutting down...")
					return
				}
			}
		}
	}
}
