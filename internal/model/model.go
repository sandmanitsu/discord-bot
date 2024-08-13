package model

import (
	"fmt"
	"sc-bot/internal/config"
	"sc-bot/internal/disk"
	"sc-bot/internal/messages"
	"strings"

	"github.com/bwmarrin/discordgo"
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

func Play(s string) string {
	// service, err := disk.GetService()
	// if err != nil {
	// 	fmt.Printf("Error getting service: %v", err)
	// }

	// disk.ListFilesInFolder(service, "1KaLJMxkFQ8daK39Sl8Do6jgeFDTkDlD7")
	return s
}

func GetChoices() []*discordgo.ApplicationCommandOptionChoice {
	service, err := disk.GetService()
	if err != nil {
		fmt.Printf("Error getting service: %v", err)
	}

	list := disk.ListFilesInFolder(service, "1KaLJMxkFQ8daK39Sl8Do6jgeFDTkDlD7") // Rainy Nights Of 1988
	choices := []*discordgo.ApplicationCommandOptionChoice{}

	for _, v := range list {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  formatName(v.Name),
			Value: v.Id,
		})
	}

	return choices
}

func formatName(name string) string {
	formName := strings.Split(name, "-")

	if len(formName) == 2 {
		return strings.TrimSpace(formName[1])
	}

	return "can't read a name"
}
