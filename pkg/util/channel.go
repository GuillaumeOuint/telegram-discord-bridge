package util

import (
	"os"
	"strings"

	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/types"
)

func LoadChannels() []types.ChannelMap {
	channels := os.Getenv("CHANNELS")
	DiscordChannels := make(types.ChannelMap)
	TelegramChannels := make(types.ChannelMap)
	channelsSlice := strings.Split(channels, ",")
	for _, channel := range channelsSlice {
		channelSplit := strings.Split(channel, ":")
		DiscordChannels[channelSplit[0]] = channelSplit[1]
		TelegramChannels[channelSplit[1]] = channelSplit[0]
	}
	return []types.ChannelMap{TelegramChannels, DiscordChannels}
}
