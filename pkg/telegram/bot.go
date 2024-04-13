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
	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/util"

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
				if update.Message.ReplyToMessage.MessageID != update.Message.MessageThreadID {
					message.IsReply = true
					message.ReplyTo = fmt.Sprintf("%v", update.Message.ReplyToMessage.MessageID)
				}
			}
			if update.Message.Photo != nil {
				message.Attachment = true
				if update.Message.Photo[len(update.Message.Photo)-1].FileID != "" {
					url, err := b.bot.GetFileDirectURL(update.Message.Photo[len(update.Message.Photo)-1].FileID)
					if err != nil {
						fmt.Println(err)
						break
					}
					name := update.Message.Photo[3].FileID
					// if name doesn't contain the extension, add it
					if !strings.Contains(name, ".") {
						name += ".jpg"
					}
					message.Attachments = append(message.Attachments, types.Attachement{
						URL:  url,
						Name: update.Message.Photo[len(update.Message.Photo)-1].FileID,
					})
				}
				message.Content = update.Message.Caption
			}
			if update.Message.Audio != nil {
				message.Attachment = true
				if update.Message.Audio.FileID != "" {
					url, err := b.bot.GetFileDirectURL(update.Message.Audio.FileID)
					if err != nil {
						fmt.Println(err)
						break
					}
					message.Attachments = append(message.Attachments, types.Attachement{
						URL:  url,
						Name: update.Message.Audio.FileID,
					})
				}
				message.Content = update.Message.Caption
			}
			if update.Message.Voice != nil {
				message.Attachment = true
				if update.Message.Voice.FileID != "" {
					url, err := b.bot.GetFileDirectURL(update.Message.Voice.FileID)
					if err != nil {
						fmt.Println(err)
						break
					}
					message.Attachments = append(message.Attachments, types.Attachement{
						URL:   url,
						Name:  update.Message.Voice.FileID,
						Voice: true,
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
			var ipp []interface{}
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
					Bytes: byt,
				}
				mimedetected := http.DetectContentType(byt)
				var mime util.Mime
				for _, v := range util.MimeArray {
					if mimedetected == string(v) {
						mime = v
						break
					}
				}
				if mime == "" {
					mime = util.MimeOctetStream
				}
				if mime != util.MimeOctetStream {
					extension := util.MimeExtensionMap[mime]

					currentextension := strings.Split(image.Name, ".")
					if len(currentextension) == 0 {
						image.Name = image.Name + "." + string(extension)
					} else {
						if currentextension[len(currentextension)-1] != string(extension) {
							image.Name = image.Name + "." + string(extension)
						}
					}
				}
				file.Name = image.Name
				var ip interface{}
				if mime == util.MimeImageJPEG || mime == util.MimeImagePNG || mime == util.MimeImageGIF || mime == util.MimeImageBMP || mime == util.MimeImageSVG {
					ip = tgbotapi.NewInputMediaPhoto(file)
				} else if mime == util.MimeAudioMPEG || mime == util.MimeAudioOGG || mime == util.MimeAudioWAV {
					ip = tgbotapi.NewInputMediaAudio(file)
				} else if mime == util.MimeVideoMPEG || mime == util.MimeVideoMP4 || mime == util.MimeVideoOGG || mime == util.MimeVideoQuickTime {
					ip = tgbotapi.NewInputMediaVideo(file)
				} else {
					ip = tgbotapi.NewInputMediaDocument(file)
				}
				ipp = append(ipp, ip)
			}
			msgi := tgbotapi.NewMediaGroup(chatID, ipp)
			if msg.MessageThreadID != 1 && msg.MessageThreadID != 0 {
				msgi.BaseChat.MessageThreadID = msg.MessageThreadID
			}
			resp, err := b.bot.Request(msgi)
			if err != nil {
				return err
			}
			fmt.Println(string(resp.Result))
			var mess []tgbotapi.Message
			err = json.Unmarshal(resp.Result, &mess)
			if err != nil {
				return err
			}
			message.TelegramMessageID = strconv.Itoa(mess[0].MessageID)
			editEvent := tgbotapi.EditMessageCaptionConfig{
				BaseEdit: tgbotapi.BaseEdit{
					BaseChatMessage: tgbotapi.BaseChatMessage{
						MessageID: mess[0].MessageID,
						ChatConfig: tgbotapi.ChatConfig{
							ChatID: chatID,
						},
					},
				},
				Caption: msgStr,
			}
			_, err = b.bot.Send(editEvent)
			if err != nil {
				return err
			}
			return nil
		}
		if message.Content == "" && len(message.Attachments) == 0 {
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
