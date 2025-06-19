package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/genai"
)

// GeminiClient is a client for the Gemini model.
type GeminiClient struct {
	client    *genai.Models
	modelName string
	cfg       Config
}

// NewGeminiClient creates a new Gemini client, supporting both Google AI and Vertex AI backends.
func NewGeminiClient(ctx context.Context, cfg Config, modelName string) (*GeminiClient, error) {
	var client *genai.Client
	var err error

	config := &genai.ClientConfig{}
	if apiKey := os.Getenv("GOOGLE_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	} else {
		config.Project = cfg.ProjectID
		config.Location = cfg.RegionID
		config.Backend = genai.BackendVertexAI
	}

	client, err = genai.NewClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating a genai client: %v", err)
	}

	return &GeminiClient{
		client:    client.Models,
		modelName: modelName,
		cfg:       cfg,
	}, nil
}

// GenerateContent generates content from the Gemini model.
func (c *GeminiClient) GenerateContent(ctx context.Context, w io.Writer, prompt string, parameters map[string]interface{}) error {
	var config *genai.GenerateContentConfig
	if c.cfg.ConfigFile != "" {
		modelConfig, err := os.ReadFile(c.cfg.ConfigFile)
		if err != nil {
			return fmt.Errorf("error reading model config file: %w", err)
		}
		err = json.Unmarshal(modelConfig, &config)
		if err != nil {
			return fmt.Errorf("error unmarshalling GenerationConfig from file: %w", err)
		}
		if c.cfg.LogType != "none" {
			log.Printf("config: %v", config)
		}
	}

	for result, err := range c.client.GenerateContentStream(ctx, c.modelName, genai.Text(prompt), config) {
		if err != nil {
			return err
		}
		if c.cfg.OutputType == "json" {
			rb, _ := json.MarshalIndent(result, "", "  ")
			fmt.Fprintln(w, string(rb))
		} else {
			fmt.Fprint(w, result.Text())
		}
	}

	return nil
}
