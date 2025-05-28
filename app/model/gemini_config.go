package model

type Config struct {
	GeminiPrompts []GeminiPrompt `json:"gemini_prompts"`
}

type GeminiPrompt struct {
	ID           string `json:"id"`
	Prompt       string `json:"prompt"`
	ResponseType string `json:"response_type"`
}

type Root struct {
	Config Config `json:"config"`
}
