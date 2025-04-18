package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/api/httpbody"
)

// LlamaRequest is the request to the Llama model.
type LlamaRequest struct {
	AnthropicVersion string             `json:"anthropic_version,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	Messages         []AnthropicMessage `json:"messages,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	TopP             float32            `json:"top_p,omitempty"`
	TopK             int                `json:"top_k,omitempty"`
	Temperature      float32            `json:"temperature,omitempty"`
}

// LlamaMessage is a message to the Llama model.
type LlamaMessage struct {
	Content []AnthropicContent `json:"content,omitempty"`
	Role    string             `json:"role,omitempty"`
}

// LlamaContent is the content of a message to the Llama model.
type LlamaContent struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

// AnthropicResponse is the response from the Llama model.
type LlamaResponse struct {
	ID           string         `json:"id,omitempty"`
	Type         string         `json:"type,omitempty"`
	Role         string         `json:"role,omitempty"`
	Content      []LlamaContent `json:"content,omitempty"`
	Model        string         `json:"model,omitempty"`
	StopReason   string         `json:"stop_reason,omitempty"`
	StopSequence string         `json:"stop_sequence,omitempty"`
	Usage        LlamaUsage     `json:"usage,omitempty"`
}

// TODO - Add documentation for type
type LlamaUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

// TODO - Add documentation for function
func UseLlamaModel(ctx context.Context, modelName string, cfg Config, args []string) error {
	// TODO - There are differences between this function and the matching function in palm.go, and gemini.go
	// due to when the config file contents are read.

	if cfg.LogType != "quiet" {
		log.Printf("Meta [%s]", modelName)
	}
	prompt := args[0]
	parameters := map[string]interface{}{
		//"temperature":     0.8,
		"maxTokens": 256,
		//"topP":            0.4,
		//"topK":            40,
	}
	var buf bytes.Buffer
	if err := generateContentClaude(ctx, modelName, cfg, &buf, prompt, parameters); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

// generateContentLlama generates text from prompt and configurations provided.
func generateContentLlama(ctx context.Context, modelName string, cfg Config, w io.Writer, prompt string, parameters map[string]interface{}) error {

	// Resolve unused argument
	_ = parameters

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", cfg.RegionID)

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		fmt.Fprintf(w, "unable to create prediction client: %v", err)
		return err
	}
	defer client.Close()

	// PredictRequest requires an endpoint, instances, and parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", cfg.ProjectID, cfg.RegionID, "anthropic")
	url := fmt.Sprintf("%s/%s", base, modelName)
	if cfg.LogType != "none" {
		log.Printf("url: %s", url)
	}

	// Construct an Anthropic message.
	llamaRequest := LlamaRequest{
		AnthropicVersion: "vertex-2023-10-16",
		MaxTokens:        256,
		Stream:           false,
		Messages: []AnthropicMessage{
			{
				Content: []AnthropicContent{
					{
						Text: prompt,
						Type: "text",
					},
				},
				Role: "user",
			},
		},
	}

	data, err := json.Marshal(&llamaRequest)
	if err != nil {
		return fmt.Errorf("error marshalling LlamaRequest: %v", err)
	}

	// using RawPredict
	req := &aiplatformpb.RawPredictRequest{
		Endpoint: url,
		HttpBody: &httpbody.HttpBody{
			ContentType: "application/json",
			Data:        data,
		},
	}

	resp, err := client.RawPredict(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "error in prediction: %v", err)
		return err
	}

	if cfg.OutputType == "json" {
		fmt.Fprintln(w, string(resp.Data))
	} else {
		var r LlamaResponse
		_ = json.Unmarshal(resp.Data, &r)
		fmt.Fprintf(w, "%v", r.Content[0].Text)

	}

	return nil
}
