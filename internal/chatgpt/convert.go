package chatgpt

import (
	internal_types "freechatgpt/internal/types"
)

func ConvertAPIRequest(api_request internal_types.APIRequest) internal_types.ChatGPTRequest {
	chatgpt_request := internal_types.NewChatGPTRequest()
	for _, api_message := range api_request.Messages {
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	return chatgpt_request
}
