package model

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

// UseGeminiModel calls Gemini's generate content method
func UseGeminiModel(ctx context.Context, modelName string, cfg Config, args []string) error {
	log.Printf("Gemini [%s]", modelName)
	prompt := genai.Text(args[0])
	var buf bytes.Buffer
	if err := GenerateContentGemini(ctx, modelName, cfg, &buf, []genai.Part{prompt}); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

// GenerateContentGemini calls Gemini's generate content method
func GenerateContentGemini(ctx context.Context, modelName string, cfg Config, w io.Writer, parts []genai.Part) error {

	client, err := genai.NewClient(ctx, cfg.ProjectID, cfg.RegionID)

	if err != nil {
		return fmt.Errorf("error creating a client: %v", err)
	}
	gemini := client.GenerativeModel(modelName)

	if cfg.ConfigFile != "" {
		modelConfig, err := os.ReadFile(cfg.ConfigFile)
		if err != nil {
			return fmt.Errorf("error reading model config file: %w", err)
		}
		var config genai.GenerationConfig
		err = json.Unmarshal(modelConfig, &config)
		if err != nil {
			return fmt.Errorf("error unmarshalling GenerationConfig from file: %w", err)
		}
		gemini.GenerationConfig = config
		if cfg.LogType != "none" {
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
	if cfg.OutputType == "json" {
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
