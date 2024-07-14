package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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
			} else if msg.Document != nil {
				handleDocument(bot, msg)
			} else if msg.IsCommand() == false {
				message := tgbotapi.NewMessage(msg.Chat.ID, "It is not subtitle file. Try again")
				bot.Send(message)
			}
		}
	}
}

func handleStartCmd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	txt := "Welcome. Send me subtitle file you want to timeshift."
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	bot.Send(msg)
}

func handleDocument(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	fileID := message.Document.FileID
	fileName := message.Document.FileName
	fileSize := message.Document.FileSize
	fileMime := message.Document.MimeType

	if fileMime == "application/x-subrip" || (fileMime == "text/plain" && strings.HasSuffix(fileName, ".srt")) {

		log.Printf("Received file: %s (ID: %s, Size: %v)", fileName, fileID, fileSize)

		// Get the file
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Printf("Error getting file URL: %s", err)
			return
		}

		fileURL := "https://api.telegram.org/file/bot" + bot.Token + "/" + file.FilePath

		err = downloadFile(fileName, fileURL)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("file downloaded successfully")
	} else {
		txt := "It is not subtitle file. Try again"
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)
		bot.Send(msg)
	}
}

func downloadFile(fileName, fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
