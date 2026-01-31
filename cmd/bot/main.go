package main

import (
	"log"
	bot "study-bot/internal/bot"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot.Run()
}
