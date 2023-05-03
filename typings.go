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

type APIRequest struct {
	Model            string             `json:"model"`
	Messages         []Message          `json:"messages"`
	Name             string             `json:"name,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	N                int                `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	Stop             interface{}        `json:"stop,omitempty"`
	MaxTokens        *int               `json:"max_tokens,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
}
