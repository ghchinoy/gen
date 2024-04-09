package cmd

// A Model sends prompts to a specific GenAI model using an Endpoint location, where the model is enabled and billed
type Model struct {
	Prompt func(projectID string, region string, modelName string, args []string) error
	// endpoint Endpoint
	mFamily string
	mType   string
	mName   string
}

var Models map[string]Model = map[string]Model{
	"gemini-1.0-pro-001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "gemini-1.0-pro-001",
	},
	"gemini-1.0-ultra-001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "gemini-1.0-ultra-001",
	},
	"gemini-1.0-pro-vision-001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "gemini-1.0-pro-vision-001",
	},
	"gemini-1.0-ultra-vision-001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "gemini-1.0-ultra-vision-001",
	},
	"gemini-1.5-pro-preview-0215": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "gemini-1.5-pro-preview-0215",
	},
	"text-bison": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "text-bison",
	},
	"text-bison@001": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "text-bison@001",
	},
	"text-bison@002": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "text-bison@002",
	},
	"text-unicorn@001": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "text-unicorn@001",
	},
	"medlm-medium": {
		Prompt:  usePaLMModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "medlm-medium",
	},
	"medlm-large": {
		Prompt:  usePaLMModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "medlm-large",
	},
	"medpalm2@preview": {
		Prompt:  usePaLMModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "medpalm2@preview",
	},
	"code-bison": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-bison",
	},
	"code-bison@001": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-bison@001",
	},
	"code-bison@002": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-bison@002",
	},
	"code-bison-32k": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-bison-32k",
	},
	"code-bison-32k@002": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-bison-32k@002",
	},
	"code-gecko": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-gecko",
	},
	"code-gecko@001": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-gecko@001",
	},
	"code-gecko@002": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "code-gecko@002",
	},
	"claude-3-haiku@20240307": {
		Prompt:  useClaudeModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "claude-3-haiku@20240307",
	},
}
