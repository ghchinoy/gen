package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"cloud.google.com/go/vertexai/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	modelName string
)

func init() {
	rootCmd.AddCommand(promptCmd)

	//promptCmd.AddCommand(generateContentCmd)

	promptCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-1.0-pro", "model name")

	//flag.StringVar(&modelName, "model", "gemini-1.0-pro", "generative model to use")
	//flag.StringVar(&region, "region", "us-central1", "region to use")
	//flag.Parse()
}

var promptCmd = &cobra.Command{
	Use:     "prompt",
	Aliases: []string{"p"},
	Short:   "Prompt a model",
	Long:    `Provide prompt parts to a model to generate content`,
	Run:     generateContentForModel,
}

// generateContentForModel prompts a model to generate content based on the provided prompt.
func generateContentForModel(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("please provide prompt")
		os.Exit(1)
	}
	log.Printf("project / region: %s / %s", projectID, region)
	log.Printf("model: %s", modelName)
	log.Printf("prompt: %s", args)

	// TODO better as a switch guarded by model list
	var err error
	if strings.HasPrefix(modelName, "gemini") {
		err = useGeminiModel(projectID, region, modelName, args)
	} else if strings.Contains(modelName, "bison") {
		err = usePaLMModel(projectID, region, modelName, args)
	} else if strings.HasPrefix(modelName, "medlm-") || strings.HasPrefix(modelName, "medpalm") {
		err = usePaLMModel(projectID, region, modelName, args)
	} else if strings.HasPrefix(modelName, "claude") {
		err = useClaudeModel(projectID, region, modelName, args)
	} else {
		err = fmt.Errorf("model '%s' is not supported", modelName)
	}
	if err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
}

// useGeminiModel calls Gemini's generate content method
func useGeminiModel(projectID string, region string, modelName string, args []string) error {
	log.Print("using Gemini")
	prompt := genai.Text(args[0])
	var buf bytes.Buffer
	if err := generateContentGemini(&buf, projectID, region, modelName, []genai.Part{prompt}); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

// usePaLMModel calls PaLM's generate content method
func usePaLMModel(projectID string, region string, modelName string, args []string) error {
	log.Print("using PaLM2")
	prompt := args[0]
	parameters := map[string]interface{}{
		"temperature":     0.8,
		"maxOutputTokens": 256,
		"topP":            0.4,
		"topK":            40,
	}
	var buf bytes.Buffer
	if err := generateContentPaLM(&buf, prompt, projectID, region, "google", modelName, parameters); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	log.Printf("generated content: %s", buf.String())
	return nil
}

func useClaudeModel(projectID string, region string, modelName string, args []string) error {
	log.Print("using Anthropic")
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
	log.Printf("generated content: %s", buf.String())
	return nil
}

// generateContentGemini calls Gemini's generate content method
func generateContentGemini(w io.Writer, projectID string, region string, modelName string, parts []genai.Part) error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, region)
	if err != nil {
		return fmt.Errorf("error creating a client: %v", err)
	}
	gemini := client.GenerativeModel(modelName)

	resp, err := gemini.GenerateContent(ctx, parts...)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	if outputtype == "json" {
		rb, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Fprintln(w, string(rb))
	} else {
		if len(resp.Candidates) > 0 {
			var all []string
			for _, v := range resp.Candidates[0].Content.Parts {
				all = append(all, fmt.Sprintf("%s", v))
			}
			fmt.Fprintf(w, "%s", strings.Join(all, " "))
		} else {
			log.Printf("Candidate length 0")
		}
	}
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
	log.Printf("url: %s", url)

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

	fmt.Fprintf(w, "text-prediction response: %v", resp.Predictions[0])
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
	log.Printf("url: %s", url)

	// Construct an Anthropic message.
	// []messages{ role, []content { type, text }, max_tokens }
	message, err := structpb.NewValue(map[string]interface{}{
		"role": "user",
		"content": []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": prompt,
			}},
	})
	if err != nil {
		log.Printf("can't create message from content: %v", err)
		return err
	}
	log.Printf("message: %v", message.GetStructValue().AsMap())

	messages, err := structpb.NewList([]interface{}{
		message.GetStructValue().AsMap(),
	})
	if err != nil {
		log.Printf("can't create messages from message: %v", err)
		return err
	}
	log.Printf("messages: %v", messages)

	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"anthropic_version": "vertex-2023-10-16",
		"max_tokens":        parameters["maxTokens"].(int),
		"stream":            false,
		"messages":          messages.AsSlice(),
	})
	if err != nil {
		log.Printf("can't create promptValue: %v", err)
		return err
	}
	jsonbytes, _ := protojson.Marshal(promptValue)
	log.Printf("jsonbytes: %s", jsonbytes)

	/*
		messagesValue, err := structpb.NewValue(map[string]interface{}{
			"messages": messages.AsSlice(),
		})
		if err != nil {
			log.Printf("can't create messagesValue: %v", err)

			return err
		}
	*/

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:  url,
		Instances: []*structpb.Value{message},
		//Parameters: parametersValue,
	}

	// PredictResponse: receive the response from the model
	resp, err := client.Predict(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "error in prediction: %v", err)
		return err
	}

	fmt.Fprintf(w, "text-prediction response: %v", resp.Predictions[0])
	return nil
}
