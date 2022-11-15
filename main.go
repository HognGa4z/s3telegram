package main

import (
	"fmt"
	"log"
	"strings"
	"s3telegram/config"
	"s3telegram/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	var conf config.BotConfig
	err := config.Read(&conf)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if !util.IsRawPhoto(update.Message) {
				continue
			}
			s3_url, err := util.ProcessDocument(update.Message.Document, bot, &conf)
			if err != nil {
				log.Fatal(err)
				continue
			}
			log.Printf("s3 location [%s]", s3_url)
			s3_url = strings.Replace(s3_url, conf.S3Host, conf.CloudFrontHost, -1)
			log.Printf("cloud front location [%s]", s3_url)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("s3 location\n[%s]", s3_url))
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
