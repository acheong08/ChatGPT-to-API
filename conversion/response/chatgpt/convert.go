package chatgpt

import (
	chatgpt_types "freechatgpt/typings/chatgpt"
	official_types "freechatgpt/typings/official"
	"strings"
)

var Previous_text string

func ConvertToString(chatgpt_response *chatgpt_types.ChatGPTResponse) string {
	translated_response := official_types.NewChatCompletionChunk(strings.ReplaceAll(chatgpt_response.Message.Content.Parts[0], Previous_text, ""))
	Previous_text = chatgpt_response.Message.Content.Parts[0]
	return "data:" + translated_response.String() + "\n\n"

}
