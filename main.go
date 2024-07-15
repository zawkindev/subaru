package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"subaru/utility"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	shiftingMode := false
	fileName := ""
	for update := range updates {
		msg := update.Message
		if msg != nil {
			if msg.IsCommand() {
				handleStartCmd(bot, msg)
			} else if msg.Document != nil {
				err := handleDocument(bot, msg)
				if err != nil {
					log.Println(err)
					continue
				}
				message := tgbotapi.NewMessage(msg.Chat.ID, "Ok, now send the time to timeshift in seconds.\ne.g.\n1000 -> for shifting 1s\n-1000 -> for shifting 1s back.")
				bot.Send(message)
				shiftingMode = true
				fileName = msg.Document.FileName

			} else if msg.IsCommand() == false {
				if shiftingMode {
					seconds, err := strconv.Atoi(strings.TrimSpace(msg.Text))
					if err != nil {
						log.Printf("Error converting string to int: %v", err)
						message := tgbotapi.NewMessage(msg.Chat.ID, "Invalid number format. Please send an integer.")
						bot.Send(message)
						continue
					}

					utility.TimeShift(fileName, int64(seconds)) // error is HERE!!!!!!!!

					if err = sendFile(bot, msg, fileName); err != nil {
						log.Fatal(err)
					}
					shiftingMode = false

					err = os.Remove(fileName)
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}
					fmt.Println("File removed successfully")

				} else {
					message := tgbotapi.NewMessage(msg.Chat.ID, "It is not subtitle file. Try again")
					shiftingMode = false
					bot.Send(message)
				}
			}
		}
	}
}

func handleStartCmd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	txt := "Welcome. Send me subtitle file you want to timeshift."
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	bot.Send(msg)
}

func handleDocument(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	fileID := message.Document.FileID
	fileName := message.Document.FileName
	fileSize := message.Document.FileSize
	fileMime := message.Document.MimeType

	if fileMime == "application/x-subrip" || (fileMime == "text/plain" && strings.HasSuffix(fileName, ".srt")) {

		log.Printf("Received file: %s (ID: %s, Size: %v)", fileName, fileID, fileSize)

		// Get the file
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			return err
		}

		fileURL := "https://api.telegram.org/file/bot" + bot.Token + "/" + file.FilePath

		err = downloadFile(fileName, fileURL)
		if err != nil {
			return err
		}

		log.Println("file downloaded successfully")

	} else {
		txt := "It is not subtitle file. Try again"
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)
		bot.Send(msg)
		return fmt.Errorf("not srt")
	}

	return nil
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

func sendFile(bot *tgbotapi.BotAPI, message *tgbotapi.Message, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileReader{
		Name:   fileName,
		Reader: file,
	})

	_, err = bot.Send(doc)
	if err != nil {
		return err
	}

	return nil
}
