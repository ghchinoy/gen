package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/genproto/googleapis/api/httpbody"
)

// MetaClient is a client for the Meta model.
type MetaClient struct {
	client *aiplatform.PredictionClient
	cfg    Config
}

// GenerateContent generates content from the Meta model.
func (c *MetaClient) GenerateContent(ctx context.Context, w io.Writer, prompt string, parameters map[string]interface{}) error {
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", c.cfg.ProjectID, c.cfg.RegionID, "meta")
	url := fmt.Sprintf("%s/%s", base, "llama3-8b")
	if c.cfg.LogType != "none" {
		log.Printf("url: %s", url)
	}

	// Construct a Llama message.
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

	resp, err := c.client.RawPredict(ctx, req)
	if err != nil {
		return fmt.Errorf("error in prediction: %v", err)
	}

	if c.cfg.OutputType == "json" {
		fmt.Fprintln(w, string(resp.Data))
	} else {
		var r LlamaResponse
		_ = json.Unmarshal(resp.Data, &r)
		fmt.Fprintf(w, "%v", r.Content[0].Text)

	}

	return nil
}
