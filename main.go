package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

var (
	modelName string
	region    string
	projectID string
)

func init() {
	flag.StringVar(&modelName, "model", "gemini-1.0-pro", "generative model to use")
	flag.StringVar(&region, "region", "us-central1", "region to use")
	flag.Parse()
}

func main() {
	log.Printf("echo: %s", flag.Args())

	projectID = envCheck("PROJECT_ID", "")
	if projectID == "" {
		log.Fatalf("requires PROJECT_ID")
	}

	prompt := genai.Text(strings.Join(flag.Args(), " "))
	log.Printf("prompt: %s", prompt)

	var buf bytes.Buffer
	if err := generateContentGemini(&buf, projectID, region, modelName, []genai.Part{prompt}); err != nil {
		log.Fatalf("error generating content: %v", err)
	}
	log.Printf("generated content: %s", buf.String())

}

func generateContentGemini(w io.Writer, projectID string, region string, modelName string, parts []genai.Part) error {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, region)
	gemini := client.GenerativeModel(modelName)

	resp, err := gemini.GenerateContent(ctx, parts...)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	rb, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintln(w, string(rb))
	return nil
}

// envCheck checks for an environment variable, otherwise returns default
func envCheck(environmentVariable, defaultVar string) string {
	if envar, ok := os.LookupEnv(environmentVariable); !ok {
		return defaultVar
	} else if envar == "" {
		return defaultVar
	} else {
		return envar
	}
}
