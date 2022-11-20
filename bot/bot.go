package bot

import (
	"log"
	"s3telegram/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI = nil

func GetBot() *tgbotapi.BotAPI {
	if bot != nil {
		return bot
	}

	c := config.GetConfig()

	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("authorized on account %s", bot.Self.UserName)
	return bot
}
