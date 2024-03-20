package telegram

import (
	"botbridge/types"
	"fmt"
	"os"

	tgbotapi "github.com/Lakhtiste/telegram-bot-api"
)

// Bot is the main struct for the telegram bot
type Bot struct {
	bot    *tgbotapi.BotAPI
	Events chan types.Message
}

// NewBot creates a new bot
func NewBot() *Bot {
	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	return &Bot{bot: bot, Events: make(chan types.Message)}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// handle the message
	}
}

func (b *Bot) SendMessage(message types.Message) {
	//msg := tgbotapi.NewMessage(chat, message)
	//b.bot.Send(msg)
}
