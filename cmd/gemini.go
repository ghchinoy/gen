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

	"cloud.google.com/go/vertexai/genai"
)

// useGeminiModel calls Gemini's generate content method
func useGeminiModel(projectID string, region string, modelName string, args []string) error {
	log.Printf("Gemini [%s]", modelName)
	prompt := genai.Text(args[0])
	var buf bytes.Buffer
	if err := generateContentGemini(&buf, projectID, region, modelName, []genai.Part{prompt}); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
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

	if modelConfigFile != "" {
		modelConfig, err := os.ReadFile(modelConfigFile)
		if err != nil {
			return fmt.Errorf("error reading model config file: %w", err)
		}
		var config genai.GenerationConfig
		err = json.Unmarshal(modelConfig, &config)
		if err != nil {
			return fmt.Errorf("error unmarshalling GenerationConfig from file: %w", err)
		}
		gemini.GenerationConfig = config
		if logtype != "none" {
			log.Printf("config: %v", config)
		}
	}

	resp, err := gemini.GenerateContent(ctx, parts...)
	if err != nil {
		// needs more sensible parsing of error message
		if strings.Contains(err.Error(), "lookup -aiplatform.googleapis.com:") {
			log.Print("missing REGION")
		}
		if strings.Contains(err.Error(), "RESOURCE_PROJECT_INVALID") {
			log.Print("missing PROJECT_ID")
		}
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
