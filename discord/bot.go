package discord

import (
	"botbridge/types"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	bot    *discordgo.Session
	Events chan types.Message
}

func NewBot() *Bot {
	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	return &Bot{bot: dg, Events: make(chan types.Message)}
}

func (b *Bot) Start() {
	b.bot.AddHandler(b.messageCreate)
	err := b.bot.Open()
	if err != nil {
		panic(err)
	}
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// handle the message
	//s.ChannelMessageSend(m.ChannelID, m.Content)
}

func (b *Bot) SendMessage(message types.Message) {
	//msg := tgbotapi.NewMessage(chat, message)
	//b.bot.Send(msg)
}
