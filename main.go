package main

import (
	"botbridge/discord"
	"botbridge/telegram"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	telegramBot := telegram.NewBot()
	telegramBot.Start()
	discordBot := discord.NewBot()
	discordBot.Start()

	discordMessages := telegramBot.Events
	telegramMessages := telegramBot.Events

	for {
		select {
		case message := <-discordMessages:
			go telegramBot.SendMessage(message)
		case message := <-telegramMessages:
			go discordBot.SendMessage(message)
		}
	}
}
