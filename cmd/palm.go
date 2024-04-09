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
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// usePaLMModel calls PaLM's generate content method
func usePaLMModel(projectID string, region string, modelName string, args []string) error {
	if logtype != "quiet" {
		log.Printf("PaLM 2 [%s]", modelName)
	}
	prompt := args[0]

	// parameters from config file
	var parameters map[string]interface{}
	if modelConfigFile != "" {
		var err error
		parameters, err = readModelConfigFile(modelConfigFile)
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
	if logtype != "none" {
		log.Printf("config: %v", parameters)
	}

	var buf bytes.Buffer
	if err := generateContentPaLM(&buf, prompt, projectID, region, "google", modelName, parameters); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

// generateContentPaLM generates text from prompt and configurations provided.
func generateContentPaLM(w io.Writer, prompt, projectID, location, publisher, model string, parameters map[string]interface{}) error {
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

	if outputtype == "json" {
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
