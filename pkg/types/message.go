package types

import "time"

type Message struct {
	Username          string
	Content           string
	Channel           string
	DiscordMessageID  string
	TelegramMessageID string
	Date              time.Time
	IsReply           bool
	ReplyTo           string
	DiscordGuildID    string
	Attachment        bool
	Attachments       []Attachement
}

type Attachement struct {
	URL   string
	Name  string
	Voice bool
}
