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
	"time"
	_ "time/tzdata"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn := repository.Open(os.Getenv("DB_URL"))
	defer conn.Close()

	tz := os.Getenv("TZ")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Your timezone string isn't invliad %s", tz)
	}
	time.Local = loc
	log.Printf("Timezone : %s", loc.String())

	web := web.NewWeb(conn)
	bot := bot.NewBot(conn)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		web.Run(ctx)
	}()

	bot.Run(ctx)
}
