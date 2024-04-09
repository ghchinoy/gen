package cmd

// A prompter is an interface that wraps prompting a generative model at a specific endpoint location
type prompter interface {
	// prompt(ctx context.Context, args []string) error
	prompt(projectID string, region string, modelName string, args []string) error
}

// A Model sends prompts to a specific GenAI model using an Endpoint location, where the model is enabled and billed
type Model struct {
	Prompt   func(projectID string, region string, modelName string, args []string) error
	endpoint Endpoint
	mFamily  string
	mType    string
	mName    string
}

var Models map[string]Model = map[string]Model{
	"Gemini_1_0_pro_001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "Gemini_1_0_pro_001",
	},
	"Gemini_1_0_ultra_001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "Gemini_1_0_ultra_001",
	},
	"Gemini_1_0_pro_vision_001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "Gemini_1_0_pro_vision_001",
	},
	"Gemini_1_0_ultra_vision_001": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "Gemini_1_0_ultra_vision_001",
	},
	"Gemini_1_5_pro_preview_0215": {
		Prompt:  useGeminiModel,
		mFamily: "Gemini",
		mType:   "text",
		mName:   "Gemini_1_5_pro_preview_0215",
	},
	"Text_bison": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "Text_bison",
	},
	"Text_bison_001": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "Text_bison_001",
	},
	"Text_bison_002": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "Text_bison_002",
	},
	"Text_unicorn_001": {
		Prompt:  usePaLMModel,
		mFamily: "text",
		mType:   "text",
		mName:   "Text_unicorn_001",
	},
	"Medlm_medium": {
		Prompt:  usePaLMModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "Medlm_medium",
	},
	"Medlm_large": {
		Prompt:  usePaLMModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "Medlm_large",
	},
	"Medpalm2_preview": {
		Prompt:  usePaLMModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "Medpalm2_preview",
	},
	"Code_bison": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_bison",
	},
	"Code_bison_001": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_bison_001",
	},
	"Code_bison_002": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_bison_002",
	},
	"Code_bison_32k": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_bison_32k",
	},
	"Code_bison_32k_002": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_bison_32k_002",
	},
	"Code_gecko": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_gecko",
	},
	"Code_gecko_001": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_gecko_001",
	},
	"Code_gecko_002": {
		mFamily: "Embeddings",
		mType:   "Embeddings",
		mName:   "Code_gecko_002",
	},
	"Claude_3_haiku_20240307": {
		Prompt:  useClaudeModel,
		mFamily: "MultiModal",
		mType:   "MultiModal",
		mName:   "Claude_3_haiku_20240307",
	},
}
