package discord

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/db"
	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/types"
	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/util"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	bot      *discordgo.Session
	Events   chan *types.Message
	Channels types.ChannelMap
	DB       *db.DB
}

func NewBot(channels types.ChannelMap, messageDB *db.DB) *Bot {
	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences
	if err != nil {
		panic(err)
	}
	return &Bot{bot: dg, Events: make(chan *types.Message), Channels: channels, DB: messageDB}
}

func (b *Bot) Start() {
	b.bot.AddHandler(b.messageCreate)
	b.bot.AddHandler(b.memberJoinded)
	err := b.bot.Open()
	if err != nil {
		panic(err)
	}
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		panic(err)
	}
	members := guild.Members

	displayName := ""
	for _, member := range members {
		if member.User.ID == m.Author.ID {
			// Use nickname if it exists, otherwise use username
			displayName = member.Nick
			if displayName == "" {
				displayName = m.Author.GlobalName
			}
			break
		}
	}

	message := types.Message{
		Username:         displayName + " (" + m.Author.Username + ")",
		Content:          m.Content,
		Channel:          m.ChannelID,
		Date:             time.Now(),
		DiscordMessageID: m.ID,
		DiscordGuildID:   m.GuildID,
	}
	if m.ReferencedMessage != nil {
		message.IsReply = true
		message.ReplyTo = m.ReferencedMessage.ID
	}
	if len(m.Attachments) != 0 {
		message.Attachment = true
		for _, attachment := range m.Attachments {
			message.Attachments = append(message.Attachments, types.Attachement{
				URL:  attachment.URL,
				Name: attachment.Filename,
			})
		}
	}
	b.Events <- &message
}

func (b *Bot) memberJoinded(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.User.ID == s.State.User.ID {
		return
	}
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		panic(err)
	}
	members := guild.Members

	displayName := ""
	for _, member := range members {
		if member.User.ID == m.User.ID {
			// Use nickname if it exists, otherwise use username
			displayName = member.Nick
			if displayName == "" {
				displayName = m.User.GlobalName
			}
			break
		}
	}
	channelId := ""
	for i := range b.Channels {
		channelId = i
		break
	}
	message := types.Message{
		Username:         displayName + " (" + m.User.Username + ")",
		Content:          "a rejoint le groupe",
		Channel:          channelId,
		Date:             time.Now(),
		DiscordMessageID: "",
		DiscordGuildID:   m.GuildID,
	}
	b.Events <- &message
}

func (b *Bot) SendMessage(message *types.Message) error {
	discordChannel := b.Channels[message.Channel]
	if discordChannel == "" {
		for _, v := range b.Channels {
			fmt.Println(v)
			discordChannel = v
			break
		}
	}
	msgStr := message.Username + ": " + message.Content
	guildID := ""
	messageoption := &discordgo.MessageSend{
		Content: msgStr,
	}
	if message.IsReply {
		discordMessageID := 0
		for _, m := range b.DB.GetMessages() {
			if m.TelegramMessageID == message.ReplyTo {
				r, err := strconv.Atoi(m.DiscordMessageID)
				if err != nil {
					return err
				}
				discordMessageID = r
				guildID = m.DiscordGuildID
				break
			}
		}
		if discordMessageID != 0 {
			// Create a MessageReference for the message we're replying to
			messageReference := discordgo.MessageReference{}
			messageReference.MessageID = strconv.Itoa(discordMessageID)
			messageReference.ChannelID = discordChannel
			messageReference.GuildID = guildID
			// Use ChannelMessageSendReply to send the reply
			messageoption.Reference = &messageReference
		}
	}
	if message.Attachment {
		for _, imag := range message.Attachments {
			// download the image
			image, err := http.DefaultClient.Get(imag.URL)
			if err != nil {
				return err
			}
			defer image.Body.Close()
			data, err := io.ReadAll(image.Body)
			if err != nil {
				return err
			}
			mimedetected := http.DetectContentType(data)
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
			if imag.Voice {
				mime = util.MimeAudioOGG
			}
			if mime != util.MimeOctetStream {
				extension := util.MimeExtensionMap[mime]

				currentextension := strings.Split(imag.Name, ".")
				if len(currentextension) == 0 {
					imag.Name = imag.Name + "." + string(extension)
				} else {
					if currentextension[len(currentextension)-1] != string(extension) {
						imag.Name = imag.Name + "." + string(extension)
					}
				}
			}

			finimage := io.NopCloser(io.Reader(io.MultiReader(bytes.NewReader(data))))
			messageoption.Files = append(messageoption.Files, &discordgo.File{
				Name:   imag.Name,
				Reader: finimage,
			})
		}
	}
	if messageoption.Content == "" && len(messageoption.Files) == 0 {
		return errors.New("empty message")
	}
	sent, err := b.bot.ChannelMessageSendComplex(discordChannel, messageoption)
	if err != nil {
		return err
	}
	message.DiscordMessageID = sent.ID
	return nil
}
