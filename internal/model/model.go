package model

import (
	"fmt"
	"sc-bot/internal/config"
	"sc-bot/internal/disk"
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

func Play() string {
	service, err := disk.GetService()
	if err != nil {
		fmt.Printf("Error getting service: %v", err)
	}

	disk.ListFilesInFolder(service, "1KaLJMxkFQ8daK39Sl8Do6jgeFDTkDlD7")

	return "play"
}
