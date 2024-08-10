package messages

import "strings"

type MessageHistory struct {
	Messages []Messages
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func New() *MessageHistory {
	return &MessageHistory{
		Messages: make([]Messages, 0),
	}
}

func (m *MessageHistory) AppendToHistory(role, message string) {
	m.Messages = append(m.Messages, Messages{
		Role:    role,
		Content: message,
	})
}

func (m *MessageHistory) JoinMessages() string {
	var messagesParts []string
	for _, v := range m.Messages {
		messagesParts = append(messagesParts, v.Role+": "+v.Content)
	}

	return strings.Join(messagesParts, "\n")
}
