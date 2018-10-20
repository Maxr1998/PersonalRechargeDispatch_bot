package main

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Sends a message into the specified chat
func reply(bot *tgbotapi.BotAPI, m *tgbotapi.Message, message string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, message)
	bot.Send(msg)
}

func handleMessage(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	var from = m.From.ID
	var isAdmin = from == AdminID
	if m.IsCommand() {
		if !m.Chat.IsPrivate() {
			msg := tgbotapi.NewMessage(m.Chat.ID, "Please only issue commands in private chats")
			msg.ReplyToMessageID = m.MessageID
			bot.Send(msg)
			return
		}
		switch m.Command() {
		case "register":
			if isAdmin {
				log.Println("Admin is starting a registration process")
				isInRegisterMode[from] = true
				reply(bot, m, "Now forward me a message from the channel to add it to the whitelist")
			} else {
				reply(bot, m, "You have no permission to use this command")
			}
		case "start":
			if !isAdmin {
				log.Printf("User %d is attempting to register\n", from)
				isInRegisterMode[from] = true
				reply(bot, m, "Hey there! Please forward me a message from your recharge dispatch channel to register.")
			}
		case "add":
			for _, ch := range config.Channels {
				args := strings.Split(m.CommandArguments(), "\n")
				for _, arg := range args {
					if isPortalCode(arg) {
						ch.AddPortal(from, arg)
					}
				}
			}
			config.changed()
		case "save":
			config.changed()
		}
	} else if m.Chat.IsPrivate() && isInRegisterMode[from] {
		isInRegisterMode[from] = false
		if forwardChat := m.ForwardFromChat; forwardChat != nil && forwardChat.IsChannel() {
			if isAdmin {
				config.AddChannel(forwardChat.ID)
				log.Printf("Registered channel %s (%d)\n", forwardChat.Title, forwardChat.ID)
			}
			if config.GetChannel(forwardChat.ID) != nil {
				config.AddUser(forwardChat.ID, from)
				reply(bot, m, "Successfully registered")
				log.Printf("Registered user %d to channel %s (%d)\n", from, forwardChat.Title, forwardChat.ID)
			} else {
				reply(bot, m, "Sorry, this channel isn't whitelisted in the bot")
			}
		} else {
			reply(bot, m, "Registration failed, please forward me a message from a (dispatch) channel")
		}
	}
}

func handlePost(bot *tgbotapi.BotAPI, channel *Channel, text string) {
	if colonIndex := strings.IndexByte(text, ':'); colonIndex > 0 && len(text) > colonIndex+2 {
		text = text[colonIndex+2:]
	}
	var haystack string
	if rc := strings.Index(text, "Recharge: "); rc > 0 {
		haystack = text[rc+10:]
	} else if dc := strings.Index(text, "DO NOT RECHARGE: "); dc > 0 {
		haystack = text[dc+17:]
	} else {
		return
	}
	for _, v := range strings.Split(haystack, " ") {
		if isPortalCode(v) {
			log.Println(v)
			if users := channel.Portals[v]; users != nil {
				for _, user := range users {
					msg := tgbotapi.NewMessage(int64(user), text)
					bot.Send(msg)
				}
			}
		}
	}
}
