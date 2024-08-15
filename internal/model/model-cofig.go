package model

import (
	"sc-bot/internal/config"
	"sc-bot/internal/messages"
)

var (
	Token    string
	ModelURL string
	GuildId  string

	MessageHistory messages.MessageHistory
)

func init() {
	config := config.MustLoad()

	Token = config.Model.Token
	ModelURL = config.Model.ModelURL
	GuildId = config.Application.GuildID

	MessageHistory = *messages.New()
}
