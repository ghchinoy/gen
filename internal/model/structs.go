
package model

// AnthropicRequest is the request to the Anthropic model.
type AnthropicRequest struct {
	AnthropicVersion string             `json:"anthropic_version"`
	MaxTokens        int                `json:"max_tokens_to_sample"`
	Stream           bool               `json:"stream"`
	Messages         []AnthropicMessage `json:"messages"`
}

// AnthropicMessage is a message to the Anthropic model.
type AnthropicMessage struct {
	Content []AnthropicContent `json:"content"`
	Role    string             `json:"role"`
}

// AnthropicContent is the content of a message.
type AnthropicContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// AnthropicResponse is the response from the Anthropic model.
type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
}

// LlamaRequest is the request to the Llama model.
type LlamaRequest struct {
	AnthropicVersion string             `json:"anthropic_version"`
	MaxTokens        int                `json:"max_tokens_to_sample"`
	Stream           bool               `json:"stream"`
	Messages         []AnthropicMessage `json:"messages"`
}

// LlamaResponse is the response from the Llama model.
type LlamaResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
}

// PaLMResponse is the response from the PaLM model.
type PaLMResponse struct {
	Predictions []struct {
		Content string `json:"content"`
	} `json:"predictions"`
}
