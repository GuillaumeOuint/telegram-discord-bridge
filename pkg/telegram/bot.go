package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/db"
	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/types"

	tgbotapi "github.com/Lakhtiste/telegram-bot-api"
)

// Bot is the main struct for the telegram bot
type Bot struct {
	bot      *tgbotapi.BotAPI
	Events   chan *types.Message
	Channels types.ChannelMap
	DB       *db.DB
}

// NewBot creates a new bot
func NewBot(channels types.ChannelMap, messageDB *db.DB) *Bot {
	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	return &Bot{bot: bot, Events: make(chan *types.Message), Channels: channels, DB: messageDB}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if !update.Message.Chat.IsForum {
				continue
			}

			if update.Message.MessageThreadID == 0 {
				update.Message.MessageThreadID = 1
			}
			message := types.Message{
				Username:          update.Message.From.FirstName + " " + update.Message.From.LastName + " (" + update.Message.From.UserName + ")",
				Content:           update.Message.Text,
				Channel:           fmt.Sprintf("%v_%v", update.Message.Chat.ID, update.Message.MessageThreadID),
				TelegramMessageID: strconv.Itoa(update.Message.MessageID),
				Date:              update.Message.Time(),
			}
			if update.Message.ReplyToMessage != nil {
				message.IsReply = true
				message.ReplyTo = fmt.Sprintf("%v", update.Message.ReplyToMessage.MessageID)
			}
			if update.Message.Photo != nil {
				message.Attachment = true
				if update.Message.Photo[3].FileID != "" {
					url, err := b.bot.GetFileDirectURL(update.Message.Photo[3].FileID)
					if err != nil {
						panic(err)
					}
					message.Attachments = append(message.Attachments, types.Attachement{
						URL:  url,
						Name: update.Message.Photo[3].FileID,
					})
				}
				message.Content = update.Message.Caption
			}
			if update.Message.NewChatMembers != nil {
				message.Content = "a rejoint le groupe"
			}
			b.Events <- &message
		}
	}
}

func (b *Bot) SendMessage(message *types.Message) error {
	if telegramChannel := b.Channels[message.Channel]; telegramChannel != "" {
		msgStr := fmt.Sprintf("%s: %s", message.Username, message.Content)
		channelSplit := strings.Split(telegramChannel, "_")
		chatID, err := strconv.ParseInt(channelSplit[0], 10, 64)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(chatID, msgStr)
		i64, err := strconv.ParseInt(channelSplit[1], 10, 64)
		if err != nil {
			return err
		}
		if i64 != 1 && i64 != 0 {
			msg.MessageThreadID = int(i64)
		}
		if message.IsReply {
			telegramMessageID := 0
			for _, m := range *b.DB.Messages {
				if m.DiscordMessageID == message.ReplyTo {
					telegramMessageID, err = strconv.Atoi(m.TelegramMessageID)
					if err != nil {
						return err
					}
					break
				}
			}
			if telegramMessageID != 0 {
				msg.ReplyParameters.MessageID = telegramMessageID
			}
		}
		if message.Attachment {
			for _, image := range message.Attachments {
				// add image to msg
				img, err := http.DefaultClient.Get(image.URL)
				if err != nil {
					return err
				}
				defer img.Body.Close()
				byt, err := io.ReadAll(img.Body)
				if err != nil {
					return err
				}
				file := tgbotapi.FileBytes{
					Name:  "photo",
					Bytes: byt,
				}
				ip := tgbotapi.NewInputMediaPhoto(file)
				msgi := tgbotapi.NewMediaGroup(chatID, []interface{}{ip})
				if msg.MessageThreadID != 1 && msg.MessageThreadID != 0 {
					msgi.BaseChat.MessageThreadID = msg.MessageThreadID
				}
				resp, err := b.bot.Request(msgi)
				if err != nil {
					return err
				}

				var mess tgbotapi.Message
				err = json.Unmarshal(resp.Result, &mess)
				if err != nil {
					return err
				}
				message.TelegramMessageID = strconv.Itoa(mess.MessageID)
			}
			return nil
		}
		if message.Content == "" {
			return errors.New("empty message content")
		}
		sent, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
		message.TelegramMessageID = strconv.Itoa(sent.MessageID)
	}
	return nil
}
