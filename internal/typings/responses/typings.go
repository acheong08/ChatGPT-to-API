package responses

type ChatCompletionChunk struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Delta        Delta       `json:"delta"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func NewChatCompletionChunk(text string) ChatCompletionChunk {
	return ChatCompletionChunk{
		ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
		Object:  "chat.completion.chunk",
		Created: 0,
		Model:   "gpt-4",
		Choices: []Choices{
			{
				Index: 0,
				Delta: Delta{
					Content: text,
					Role:    "assistant",
				},
				FinishReason: nil,
			},
		},
	}
}

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}
type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Choice struct {
	Index        int         `json:"index"`
	Message      Msg         `json:"message"`
	FinishReason interface{} `json:"finish_reason"`
}
type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewChatCompletion(full_test string) ChatCompletion {
	return ChatCompletion{
		ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
		Object:  "chat.completion",
		Created: int64(0),
		Model:   "gpt-4",
		Usage: usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		Choices: []Choice{
			{
				Message: Msg{
					Content: full_test,
					Role:    "assistant",
				},
				Index:        0,
				FinishReason: "stop",
			},
		},
	}
}
