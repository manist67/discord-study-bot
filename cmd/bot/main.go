package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	bot "study-bot/internal/bot"
	"study-bot/internal/repository"
	"study-bot/internal/web"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn := repository.Open(os.Getenv("DB_URL"))
	defer conn.Close()

	web := web.NewWeb(conn)
	bot := bot.NewBot(conn)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		web.Run(ctx)
	}()

	bot.Run(ctx)
}
