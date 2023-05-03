package main

type Choice struct {
	Delta        Delta  `json:"delta"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type Data struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Delta struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type FullCompletion struct {
	ID      string     `json:"id"`
	Object  string     `json:"object"`
	Created int64      `json:"created"`
	Model   string     `json:"model"`
	Usage   Usage      `json:"usage"`
	Choices []Response `json:"choices"`
}

func NewFullCompletion(fulltext, model string) *FullCompletion {
	return &FullCompletion{
		ID:      "I am so bored",
		Object:  "chat.completion",
		Created: 0,
		Model:   model,
		Usage: Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		Choices: []Response{
			{
				Message: Message{
					Role:    "assistant",
					Content: fulltext,
				},
				FinishReason: "stop",
				Index:        0,
			},
		},
	}
}
