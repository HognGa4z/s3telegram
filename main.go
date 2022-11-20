package main

import (
	"fmt"
	"log"
	"s3telegram/bot"
	"s3telegram/config"
	"s3telegram/util"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	conf := config.GetConfig()
	bot := bot.GetBot()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = conf.BotTimeout
	updates := bot.GetUpdatesChan(u)

	go util.TmpFileClear()

	for update := range updates {
		if update.Message != nil { // If we got a message
			go func() {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				if !util.IsRawPhoto(update.Message) {
					return
				}
				s3_url, err := util.ProcessUploadS3(update.Message.Document, bot)
				if err != nil {
					log.Printf("process upload fail, %s", err.Error())
					return
				}
				log.Printf("s3 location [%s]", s3_url)
				s3_url = strings.Replace(s3_url, conf.S3Host, conf.CloudFrontHost, -1)
				log.Printf("cloud front location [%s]", s3_url)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("s3 location\n[%s]", s3_url))
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}()
		}
	}
}
