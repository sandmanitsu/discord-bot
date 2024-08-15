package model

func Dialog(newMessage string) string {
	// todo проверить использует ли бот историю диалога
	// MessageHistory.AppendToHistory("user", newMessage)

	return Request(newMessage)
}
