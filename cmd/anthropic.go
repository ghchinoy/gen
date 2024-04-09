package cmd

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

func useClaudeModel(projectID string, region string, modelName string, args []string) error {
	if logtype != "quiet" {
		log.Printf("Anthropic [%s]", modelName)
	}
	prompt := args[0]
	parameters := map[string]interface{}{
		//"temperature":     0.8,
		"maxTokens": 256,
		//"topP":            0.4,
		//"topK":            40,
	}
	var buf bytes.Buffer
	if err := generateContentClaude(&buf, prompt, projectID, region, "anthropic", modelName, parameters); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

// generateContentClaude generates text from prompt and configurations provided.
func generateContentClaude(w io.Writer, prompt, projectID, location, publisher, model string, parameters map[string]interface{}) error {
	ctx := context.Background()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		fmt.Fprintf(w, "unable to create prediction client: %v", err)
		return err
	}
	defer client.Close()

	// PredictRequest requires an endpoint, instances, and parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", projectID, location, publisher)
	url := fmt.Sprintf("%s/%s", base, model)
	if logtype != "none" {
		log.Printf("url: %s", url)
	}

	// Construct an Anthropic message.
	claudeRequest := AnthropicRequest{
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

	data, err := json.Marshal(&claudeRequest)
	if err != nil {
		return fmt.Errorf("error marshalling ClaudeRequest: %v", err)
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

	if outputtype == "json" {
		fmt.Fprintln(w, string(resp.Data))
	} else {
		var r AnthropicResponse
		_ = json.Unmarshal(resp.Data, &r)
		fmt.Fprintf(w, "%v", r.Content[0].Text)

	}

	return nil
}