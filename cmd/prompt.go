package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	// "google.golang.org/protobuf/encoding/protojson"
	// "google.golang.org/protobuf/types/known/structpb"
)

var (
	modelName       string
	modelConfigFile string
	//modelConfig     map[string]interface{}
)

func init() {
	rootCmd.AddCommand(promptCmd)

	//promptCmd.AddCommand(generateContentCmd)

	promptCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-1.0-pro", "model name")
	promptCmd.PersistentFlags().StringVarP(&modelConfigFile, "config", "c", "", "model parameters")

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
	if logtype != "none" {
		log.Printf("model: %s", modelName)
		log.Printf("prompt: %s", args)
	}

	// TODO better as a switch guarded by model list
	var err error
	if strings.HasPrefix(modelName, "gemini") {
		err = useGeminiModel(projectID, region, modelName, args)
	} else if strings.Contains(modelName, "text-bison") || strings.Contains(modelName, "text-unicorn") {
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

// // useGeminiModel calls Gemini's generate content method
// func useGeminiModel(projectID string, region string, modelName string, args []string) error {
// 	log.Printf("Gemini [%s]", modelName)
// 	prompt := genai.Text(args[0])
// 	var buf bytes.Buffer
// 	if err := generateContentGemini(&buf, projectID, region, modelName, []genai.Part{prompt}); err != nil {
// 		log.Printf("error generating content: %v", err)
// 		os.Exit(1)
// 	}
// 	fmt.Printf("%s\n", buf.String())
// 	return nil
// }

// // usePaLMModel calls PaLM's generate content method
// func usePaLMModel(projectID string, region string, modelName string, args []string) error {
// 	if logtype != "quiet" {
// 		log.Printf("PaLM 2 [%s]", modelName)
// 	}
// 	prompt := args[0]

// 	// parameters from config file
// 	var parameters map[string]interface{}
// 	if modelConfigFile != "" {
// 		var err error
// 		parameters, err = readModelConfigFile(modelConfigFile)
// 		if err != nil {
// 			return err
// 		}
// 	} else { // default PaLM 2 model params
// 		parameters = map[string]interface{}{
// 			"temperature":     0.6,
// 			"maxOutputTokens": 256,
// 			"topP":            0.4,
// 			"topK":            40,
// 		}
// 	}
// 	if logtype != "none" {
// 		log.Printf("config: %v", parameters)
// 	}

// 	var buf bytes.Buffer
// 	if err := generateContentPaLM(&buf, prompt, projectID, region, "google", modelName, parameters); err != nil {
// 		log.Printf("error generating content: %v", err)
// 		os.Exit(1)
// 	}
// 	fmt.Printf("%s\n", buf.String())
// 	return nil
// }

// func useClaudeModel(projectID string, region string, modelName string, args []string) error {
// 	if logtype != "quiet" {
// 		log.Printf("Anthropic [%s]", modelName)
// 	}
// 	prompt := args[0]
// 	parameters := map[string]interface{}{
// 		//"temperature":     0.8,
// 		"maxTokens": 256,
// 		//"topP":            0.4,
// 		//"topK":            40,
// 	}
// 	var buf bytes.Buffer
// 	if err := generateContentClaude(&buf, prompt, projectID, region, "anthropic", modelName, parameters); err != nil {
// 		log.Printf("error generating content: %v", err)
// 		os.Exit(1)
// 	}
// 	fmt.Printf("%s\n", buf.String())
// 	return nil
// }

// // generateContentGemini calls Gemini's generate content method
// func generateContentGemini(w io.Writer, projectID string, region string, modelName string, parts []genai.Part) error {
// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, projectID, region)
// 	if err != nil {
// 		return fmt.Errorf("error creating a client: %v", err)
// 	}
// 	gemini := client.GenerativeModel(modelName)

// 	if modelConfigFile != "" {
// 		modelConfig, err := os.ReadFile(modelConfigFile)
// 		if err != nil {
// 			return fmt.Errorf("error reading model config file: %w", err)
// 		}
// 		var config genai.GenerationConfig
// 		err = json.Unmarshal(modelConfig, &config)
// 		if err != nil {
// 			return fmt.Errorf("error unmarshalling GenerationConfig from file: %w", err)
// 		}
// 		gemini.GenerationConfig = config
// 		if logtype != "none" {
// 			log.Printf("config: %v", config)
// 		}
// 	}

// 	resp, err := gemini.GenerateContent(ctx, parts...)
// 	if err != nil {
// 		// needs more sensible parsing of error message
// 		if strings.Contains(err.Error(), "lookup -aiplatform.googleapis.com:") {
// 			log.Print("missing REGION")
// 		}
// 		if strings.Contains(err.Error(), "RESOURCE_PROJECT_INVALID") {
// 			log.Print("missing PROJECT_ID")
// 		}
// 		return fmt.Errorf("error generating content: %w", err)
// 	}
// 	if outputtype == "json" {
// 		rb, _ := json.MarshalIndent(resp, "", "  ")
// 		fmt.Fprintln(w, string(rb))
// 	} else {
// 		if len(resp.Candidates) > 0 {
// 			var all []string
// 			for _, v := range resp.Candidates[0].Content.Parts {
// 				all = append(all, fmt.Sprintf("%s", v))
// 			}
// 			fmt.Fprintf(w, "%s", strings.Join(all, " "))
// 		} else {
// 			log.Printf("Candidate length 0")
// 		}
// 	}
// 	return nil
// }

// // generateContentPaLM generates text from prompt and configurations provided.
// func generateContentPaLM(w io.Writer, prompt, projectID, location, publisher, model string, parameters map[string]interface{}) error {
// 	ctx := context.Background()

// 	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

// 	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
// 	if err != nil {
// 		fmt.Fprintf(w, "unable to create prediction client: %v", err)
// 		return err
// 	}
// 	defer client.Close()

// 	// PredictRequest requires an endpoint, instances, and parameters
// 	// Endpoint
// 	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", projectID, location, publisher)
// 	url := fmt.Sprintf("%s/%s", base, model)
// 	if logtype != "none" {
// 		log.Printf("url: %s", url)
// 	}
// 	// Instances: the prompt to use with the text model
// 	promptValue, err := structpb.NewValue(map[string]interface{}{
// 		"prompt": prompt,
// 	})
// 	if err != nil {
// 		fmt.Fprintf(w, "unable to convert prompt to Value: %v", err)
// 		return err
// 	}

// 	// Parameters: the model configuration parameters
// 	parametersValue, err := structpb.NewValue(parameters)
// 	if err != nil {
// 		fmt.Fprintf(w, "unable to convert parameters to Value: %v", err)
// 		return err
// 	}

// 	// PredictRequest: create the model prediction request
// 	req := &aiplatformpb.PredictRequest{
// 		Endpoint:   url,
// 		Instances:  []*structpb.Value{promptValue},
// 		Parameters: parametersValue,
// 	}

// 	// PredictResponse: receive the response from the model
// 	resp, err := client.Predict(ctx, req)
// 	if err != nil {
// 		fmt.Fprintf(w, "error in prediction: %v", err)
// 		return err
// 	}

// 	if outputtype == "json" {
// 		rb, _ := json.MarshalIndent(resp, "", "  ")
// 		fmt.Fprintln(w, string(rb))
// 	} else {
// 		if len(resp.Predictions) > 0 {
// 			var r PaLMResponse
// 			structbytes, _ := protojson.Marshal(resp)
// 			err := json.Unmarshal(structbytes, &r)
// 			if err != nil {
// 				return fmt.Errorf("unable to convert to struct: %v", err)
// 			}
// 			fmt.Fprintf(w, "%v", r.Predictions[0].Content)
// 		}
// 	}
// 	return nil

// }

// // generateContentClaude generates text from prompt and configurations provided.
// func generateContentClaude(w io.Writer, prompt, projectID, location, publisher, model string, parameters map[string]interface{}) error {
// 	ctx := context.Background()

// 	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

// 	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
// 	if err != nil {
// 		fmt.Fprintf(w, "unable to create prediction client: %v", err)
// 		return err
// 	}
// 	defer client.Close()

// 	// PredictRequest requires an endpoint, instances, and parameters
// 	// Endpoint
// 	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", projectID, location, publisher)
// 	url := fmt.Sprintf("%s/%s", base, model)
// 	if logtype != "none" {
// 		log.Printf("url: %s", url)
// 	}

// 	// Construct an Anthropic message.
// 	claudeRequest := AnthropicRequest{
// 		AnthropicVersion: "vertex-2023-10-16",
// 		MaxTokens:        256,
// 		Stream:           false,
// 		Messages: []AnthropicMessage{
// 			{
// 				Content: []AnthropicContent{
// 					{
// 						Text: prompt,
// 						Type: "text",
// 					},
// 				},
// 				Role: "user",
// 			},
// 		},
// 	}

// 	data, err := json.Marshal(&claudeRequest)
// 	if err != nil {
// 		return fmt.Errorf("error marshalling ClaudeRequest: %v", err)
// 	}

// 	// using RawPredict
// 	req := &aiplatformpb.RawPredictRequest{
// 		Endpoint: url,
// 		HttpBody: &httpbody.HttpBody{
// 			ContentType: "application/json",
// 			Data:        data,
// 		},
// 	}

// 	resp, err := client.RawPredict(ctx, req)
// 	if err != nil {
// 		fmt.Fprintf(w, "error in prediction: %v", err)
// 		return err
// 	}

// 	if outputtype == "json" {
// 		fmt.Fprintln(w, string(resp.Data))
// 	} else {
// 		var r AnthropicResponse
// 		_ = json.Unmarshal(resp.Data, &r)
// 		fmt.Fprintf(w, "%v", r.Content[0].Text)

// 	}

// 	return nil
// }

// readModelConfigFile reads the model configuration file (JSON text file)
func readModelConfigFile(configFile string) (map[string]interface{}, error) {
	var config map[string]interface{}
	data, err := os.ReadFile(modelConfigFile)
	if err != nil {
		return config, fmt.Errorf("error reading model config: %v", err)

	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling model config: %v", err)
	}
	return config, nil
}

// AnthropicRequest is the request to the Claude model.
type AnthropicRequest struct {
	AnthropicVersion string             `json:"anthropic_version,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	Messages         []AnthropicMessage `json:"messages,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	TopP             float32            `json:"top_p,omitempty"`
	TopK             int                `json:"top_k,omitempty"`
	Temperature      float32            `json:"temperature,omitempty"`
}

// AnthropicMessage is a message to the Claude model.
type AnthropicMessage struct {
	Content []AnthropicContent `json:"content,omitempty"`
	Role    string             `json:"role,omitempty"`
}

// AnthropicContent is the content of a message to the Claude model.
type AnthropicContent struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

// AnthropicResponse is the response from the Claude model.
type AnthropicResponse struct {
	ID           string             `json:"id,omitempty"`
	Type         string             `json:"type,omitempty"`
	Role         string             `json:"role,omitempty"`
	Content      []AnthropicContent `json:"content,omitempty"`
	Model        string             `json:"model,omitempty"`
	StopReason   string             `json:"stop_reason,omitempty"`
	StopSequence string             `json:"stop_sequence,omitempty"`
	Usage        AnthropicUsage     `json:"usage,omitempty"`
}

type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

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
