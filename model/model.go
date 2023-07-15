package model

type ReqBody struct {
	Model            string  `json:"model"`
	Messages         []Roles `json:"messages"`
	Temperature      float64 `json:"temperature,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	N                int     `json:"n,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
}

type Roles struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptConfig struct {
	Temperature      float64 `json:"temperature,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	N                int     `json:"n,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
}

type ResBody struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int64      `json:"index"`
	Message      ResMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type ResMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type Env struct {
	GptToken string
}

type QnaBody struct {
	Q string
	A string
}
