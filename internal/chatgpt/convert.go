package chatgpt

import (
	typings "freegpt4/internal/typings"
)

func ConvertAPIRequest(api_request typings.APIRequest) typings.ChatMessage {
	chatgpt_request := typings.NewChatMessage()
	chatgpt_request.Meta.Content.InternetAccess = api_request.Internet
	for _, api_message := range api_request.Messages {
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	// Remove the last message
	chatgpt_request.Rollback()
	// Add the last message as a part
	chatgpt_request.AddPart(api_request.Messages[len(api_request.Messages)-1].Role, api_request.Messages[len(api_request.Messages)-1].Content)
	return chatgpt_request
}
