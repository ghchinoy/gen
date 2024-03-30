package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/vertexai/genai"
	"github.com/spf13/cobra"
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

func generateContentForModel(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("please provide prompt")
		os.Exit(1)
	}
	//log.Printf("project / region: %s / %s", projectID, region)
	log.Printf("model: %s", modelName)
	log.Printf("prompt: %s", args)

	prompt := genai.Text(args[0])
	var buf bytes.Buffer
	if err := generateContentGemini(&buf, projectID, region, modelName, []genai.Part{prompt}); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	log.Printf("generated content: %s", buf.String())
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
