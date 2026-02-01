package main

import (
	"log"
	"os"
	bot "study-bot/internal/bot"
	"study-bot/internal/repository"
	"study-bot/internal/web"

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

	go func() {
		web.Run()
	}()

	bot.Run()
}
