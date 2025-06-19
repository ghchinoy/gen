package model

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"google.golang.org/api/option"
)

// ModelClient is an interface for interacting with a generative AI model.
type ModelClient interface {
	// GenerateContent sends a prompt to the model and returns the generated content.
	GenerateContent(ctx context.Context, w io.Writer, prompt string, parameters map[string]interface{}) error
}

// NewClient creates a new model client based on the model name.
func NewClient(ctx context.Context, cfg Config, modelName string) (ModelClient, error) {
	if cfg.ProjectID == "" {
		cfg.ProjectID = os.Getenv("GEN_PROJECT_ID")
	}
	if cfg.RegionID == "" {
		cfg.RegionID = os.Getenv("GEN_REGION")
	}

	if strings.HasPrefix(modelName, "gemini") {
		return NewGeminiClient(ctx, cfg, modelName)
	}

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", cfg.RegionID)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return nil, fmt.Errorf("unable to create prediction client: %v", err)
	}

	if strings.HasPrefix(modelName, "text-bison") {
		return &PaLMClient{client: client, cfg: cfg}, nil
	} else if strings.HasPrefix(modelName, "claude") {
		return &AnthropicClient{client: client, cfg: cfg}, nil
	} else if strings.HasPrefix(modelName, "llama") {
		return &MetaClient{client: client, cfg: cfg}, nil
	}
	return nil, fmt.Errorf("unknown model: %s", modelName)
}
