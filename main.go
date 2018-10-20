package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var config *Config
var isInRegisterMode = make(map[int]bool)

func main() {
	// Load config
	config = LoadConfig()

	// Register bot with API
	bot, err := tgbotapi.NewBotAPI(APIKey)
	if err != nil {
		log.Panic(err)
		return
	}

	// Create update channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// Process updates from the Telegram servers
	for update := range updates {
		if post := update.ChannelPost; post != nil {
			handlePost(bot, config.GetChannel(post.Chat.ID), post.Text)
			continue
		}
		if msg := update.Message; msg != nil {
			handleMessage(bot, msg)
		}
	}
}
