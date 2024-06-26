package main

import (
	"fmt"

	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/db"
	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/discord"
	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/telegram"
	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/util"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	channels := util.LoadChannels()

	messageDB := db.NewDB()
	telegramBot := telegram.NewBot(channels[0], messageDB)
	go telegramBot.Start()
	discordBot := discord.NewBot(channels[1], messageDB)
	go discordBot.Start()

	discordMessages := discordBot.Events
	telegramMessages := telegramBot.Events

	for {
		select {
		case message := <-discordMessages:
			fmt.Printf("Received discord message: %#v\n", message)
			messageDB.AddMessage(message)
			go func() {
				err := telegramBot.SendMessage(message)
				if err != nil {
					fmt.Printf("Error sending message to telegram: %v\n", err)
				}
			}()
		case message := <-telegramMessages:
			fmt.Printf("Received telegram message: %#v\n", message)
			messageDB.AddMessage(message)
			go func() {
				err := discordBot.SendMessage(message)
				if err != nil {
					fmt.Printf("Error sending message to discord: %v\n", err)
				}
			}()
		}
	}
}
