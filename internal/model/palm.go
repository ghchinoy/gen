package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// PaLMClient is a client for the PaLM model.
type PaLMClient struct {
	client *aiplatform.PredictionClient
	cfg    Config
}

// GenerateContent generates content from the PaLM model.
func (c *PaLMClient) GenerateContent(ctx context.Context, w io.Writer, prompt string, parameters map[string]interface{}) error {
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", c.cfg.ProjectID, c.cfg.RegionID, "google")
	url := fmt.Sprintf("%s/%s", base, "text-bison")
	if c.cfg.LogType != "none" {
		log.Printf("url: %s", url)
	}
	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})
	if err != nil {
		return fmt.Errorf("unable to convert prompt to Value: %v", err)
	}

	// Parameters: the model configuration parameters
	parametersValue, err := structpb.NewValue(parameters)
	if err != nil {
		return fmt.Errorf("unable to convert parameters to Value: %v", err)
	}

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:   url,
		Instances:  []*structpb.Value{promptValue},
		Parameters: parametersValue,
	}

	// PredictResponse: receive the response from the model
	resp, err := c.client.Predict(ctx, req)
	if err != nil {
		return fmt.Errorf("error in prediction: %v", err)
	}

	if c.cfg.OutputType == "json" {
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
