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
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// PaLMResponse is the response from the PaLM model.
type PaLMResponse struct {
	Predictions []Prediction `json:"predictions"`
	Metadata    Metadata     `json:"metadata"`
}

type Prediction struct {
	CitationMetadata CitationMetadata `json:"citationMetadata,omitempty"`
	Content          string           `json:"content,omitempty"`
	SafetyAttributes SafetyAttributes `json:"safetyAttributes,omitempty"`
}

type CitationMetadata struct {
	Citations []interface{} `json:"citations"`
}

type SafetyAttributes struct {
	Blocked       bool           `json:"blocked,omitempty"`
	Categories    []string       `json:"categories,omitempty"`
	SafetyRatings []SafetyRating `json:"safetyRatings,omitempty"`
}

type SafetyRating struct {
	Category         string  `json:"category,omitempty"`
	ProbabilityScore float32 `json:"probabilityScore,omitempty"`
	Severity         string  `json:"severity,omitempty"`
	SeverityScore    float32 `json:"severityScore,omitempty"`
}

type Metadata struct {
	TokenMetadata TokenMetadata `json:"tokenMetadata"`
}

type TokenMetadata struct {
	InputTokenCount  TokenMetadataDetails `json:"inputTokenCount,omitempty"`
	OutputTokenCount TokenMetadataDetails `json:"outputTokenCount,omitempty"`
}

type TokenMetadataDetails struct {
	TotalBillableCharacters int `json:"totalBillableCharacters,omitempty"`
	TotalTokens             int `json:"totalTokens,omitempty"`
}

// UsePaLMModel calls PaLM's generate content method
func UsePaLMModel(ctx context.Context, modelName string, cfg Config, args []string) error {
	if cfg.LogType != "quiet" {
		log.Printf("PaLM 2 [%s]", modelName)
	}
	prompt := args[0]

	// parameters from config file
	var parameters map[string]interface{}
	if cfg.ConfigFile != "" {
		var err error
		parameters, err = cfg.ReadModelConfigFile()
		if err != nil {
			return err
		}
	} else { // default PaLM 2 model params
		parameters = map[string]interface{}{
			"temperature":     0.6,
			"maxOutputTokens": 256,
			"topP":            0.4,
			"topK":            40,
		}
	}
	if cfg.LogType != "none" {
		log.Printf("config: %v", parameters)
	}

	var buf bytes.Buffer

	if err := generateContentPaLM(ctx, modelName, cfg, &buf, prompt, parameters); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

// generateContentPaLM generates text from prompt and configurations provided.
func generateContentPaLM(ctx context.Context, modelName string, cfg Config, w io.Writer, prompt string, parameters map[string]interface{}) error {
	// TODO - There are differences between this function and the matching function in palm.go
	// due to when the config file contents are read.

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", cfg.RegionID)

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		fmt.Fprintf(w, "unable to create prediction client: %v", err)
		return err
	}
	defer client.Close()

	// PredictRequest requires an endpoint, instances, and parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", cfg.ProjectID, cfg.RegionID, "google")
	url := fmt.Sprintf("%s/%s", base, modelName)
	if cfg.LogType != "none" {
		log.Printf("url: %s", url)
	}
	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})
	if err != nil {
		fmt.Fprintf(w, "unable to convert prompt to Value: %v", err)
		return err
	}

	// Parameters: the model configuration parameters
	parametersValue, err := structpb.NewValue(parameters)
	if err != nil {
		fmt.Fprintf(w, "unable to convert parameters to Value: %v", err)
		return err
	}

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:   url,
		Instances:  []*structpb.Value{promptValue},
		Parameters: parametersValue,
	}

	// PredictResponse: receive the response from the model
	resp, err := client.Predict(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "error in prediction: %v", err)
		return err
	}

	if cfg.OutputType == "json" {
		rb, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Fprintln(w, string(rb))
	} else {
		if len(resp.Predictions) > 0 {
			var r PaLMResponse
			structbytes, _ := protojson.Marshal(resp)
			err := json.Unmarshal(structbytes, &r)
			if err != nil {
				return fmt.Errorf("unable to convert to struct: %v", err)
			}
			fmt.Fprintf(w, "%v", r.Predictions[0].Content)
		}
	}
	return nil

}
