package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		msg := update.Message
		if msg != nil {
			if msg.IsCommand() {
				handleStartCmd(bot, msg)
			}
		}
	}
}

func handleStartCmd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	txt := "Welcome. Send me subtitle file you want to timeshift."
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	bot.Send(msg)
}
