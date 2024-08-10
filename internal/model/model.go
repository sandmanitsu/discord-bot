package model

import (
	"sc-bot/internal/config"
	"sc-bot/internal/messages"
)

var (
	Token    string
	ModelURL string

	MessageHistory messages.MessageHistory
)

func init() {
	config := config.MustLoad()

	Token = config.Model.Token
	ModelURL = config.Model.ModelURL

	MessageHistory = *messages.New()
}

func Dialog(newMessage string) string {
	// todo проверить использует ли бот историю диалога
	// MessageHistory.AppendToHistory("user", newMessage)

	return Request(newMessage)
}
