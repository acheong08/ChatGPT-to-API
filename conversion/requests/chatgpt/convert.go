package chatgpt

import (
	"fmt"
	chatgpt_types "freechatgpt/typings/chatgpt"
	official_types "freechatgpt/typings/official"
	"strings"

	arkose "github.com/acheong08/funcaptcha"
)

func ConvertAPIRequest(api_request official_types.APIRequest) chatgpt_types.ChatGPTRequest {
	chatgpt_request := chatgpt_types.NewChatGPTRequest()
	if strings.HasPrefix(api_request.Model, "gpt-3.5") {
		chatgpt_request.Model = "text-davinci-002-render-sha"
	}
	if strings.HasPrefix(api_request.Model, "gpt-4") {
		token, err := arkose.GetOpenAIToken()
		if err == nil {
			chatgpt_request.ArkoseToken = token
		} else {
			fmt.Println("Error getting Arkose token: ", err)
		}
		chatgpt_request.Model = "gpt-4"
		if api_request.Model == "gpt-4-browsing" || api_request.Model == "gpt-4-mobile" || api_request.Model == "gpt-4-code-interpreter" {
			chatgpt_request.Model = api_request.Model
		}
	}
	if api_request.PluginIDs != nil {
		chatgpt_request.PluginIDs = api_request.PluginIDs
		chatgpt_request.Model = "gpt-4-plugins"
	}
	for _, api_message := range api_request.Messages {
		if api_message.Role == "system" {
			api_message.Role = "critic"
		}
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	return chatgpt_request
}
