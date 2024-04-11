package model

import (
	"context"
	"fmt"
)

// A Model sends prompts to a specific GenAI model using an Endpoint location, where the model is enabled and billed
type Model struct {
	prompt  func(ctx context.Context, modelName string, cfg Config, args []string) error
	MFamily string
	MType   string
	MName   string
}

var Models map[string]Model = map[string]Model{
	"gemini-1.0-pro": {
		prompt:  UseGeminiModel,
		MFamily: "Gemini",
		MType:   "text",
		MName:   "gemini-1.0-pro",
	},
	"gemini-1.0-pro-001": {
		prompt:  UseGeminiModel,
		MFamily: "Gemini",
		MType:   "text",
		MName:   "gemini-1.0-pro-001",
	},
	"gemini-1.0-ultra-001": {
		prompt:  UseGeminiModel,
		MFamily: "Gemini",
		MType:   "text",
		MName:   "gemini-1.0-ultra-001",
	},
	"gemini-1.0-pro-vision-001": {
		prompt:  UseGeminiModel,
		MFamily: "Gemini",
		MType:   "text",
		MName:   "gemini-1.0-pro-vision-001",
	},
	"gemini-1.0-ultra-vision-001": {
		prompt:  UseGeminiModel,
		MFamily: "Gemini",
		MType:   "text",
		MName:   "gemini-1.0-ultra-vision-001",
	},
	"gemini-1.5-pro-preview-0215": {
		prompt:  UseGeminiModel,
		MFamily: "Gemini",
		MType:   "text",
		MName:   "gemini-1.5-pro-preview-0215",
	},
	"text-bison": {
		prompt:  UsePaLMModel,
		MFamily: "text",
		MType:   "text",
		MName:   "text-bison",
	},
	"text-bison@001": {
		prompt:  UsePaLMModel,
		MFamily: "text",
		MType:   "text",
		MName:   "text-bison@001",
	},
	"text-bison@002": {
		prompt:  UsePaLMModel,
		MFamily: "text",
		MType:   "text",
		MName:   "text-bison@002",
	},
	"text-unicorn@001": {
		prompt:  UsePaLMModel,
		MFamily: "text",
		MType:   "text",
		MName:   "text-unicorn@001",
	},
	"medlm-medium": {
		prompt:  UsePaLMModel,
		MFamily: "MultiModal",
		MType:   "MultiModal",
		MName:   "medlm-medium",
	},
	"medlm-large": {
		prompt:  UsePaLMModel,
		MFamily: "MultiModal",
		MType:   "MultiModal",
		MName:   "medlm-large",
	},
	"medpalm2@preview": {
		prompt:  UsePaLMModel,
		MFamily: "MultiModal",
		MType:   "MultiModal",
		MName:   "medpalm2@preview",
	},
	"code-bison": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-bison",
	},
	"code-bison@001": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-bison@001",
	},
	"code-bison@002": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-bison@002",
	},
	"code-bison-32k": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-bison-32k",
	},
	"code-bison-32k@002": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-bison-32k@002",
	},
	"code-gecko": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-gecko",
	},
	"code-gecko@001": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-gecko@001",
	},
	"code-gecko@002": {
		MFamily: "Embeddings",
		MType:   "Embeddings",
		MName:   "code-gecko@002",
	},
	"claude-3-haiku@20240307": {
		prompt:  UseClaudeModel,
		MFamily: "MultiModal",
		MType:   "MultiModal",
		MName:   "claude-3-haiku@20240307",
	},
}

// TODO - Ideally would like to avoid this level of indirection, but suing it for the
//
//	time being to get course grained refactoring working
func (m Model) Use(ctx context.Context, cfg Config, args []string) error {
	if m.prompt != nil {
		return m.prompt(ctx, m.MName, cfg, args)
	}

	return fmt.Errorf("Model: `%s` does not currently implement the `Use` method", m.MName)

}
