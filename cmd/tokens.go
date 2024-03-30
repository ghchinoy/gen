package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/vertexai/genai"
	"github.com/spf13/cobra"
)

var promptFile string

func init() {
	rootCmd.AddCommand(tokensCmd)

	tokensCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-1.0-pro", "model name")

	tokensCmd.PersistentFlags().StringVarP(&promptFile, "file", "f", "", "prompt file")
}

var tokensCmd = &cobra.Command{
	Use:     "tokens",
	Aliases: []string{"t", "count", "tokencount", "tc"},
	Short:   "count tokens for prompt",
	Long:    `Returns the count of tokens for a provided prompt`,
	Run:     countTokensForPrompt,
}

// countTokensForPrompt is the cobra implementation of countTokens
func countTokensForPrompt(cmd *cobra.Command, args []string) {
	var prompt string
	if promptFile != "" { // read in file
		promptBytes, err := os.ReadFile(promptFile)
		if err != nil {
			log.Fatal(err)
		}
		prompt = string(promptBytes)
	} else {
		if len(args) == 0 {
			log.Fatal("requires a prompt to count tokens")
		}
		prompt = args[0]
	}

	err := countTokens(os.Stdout, string(prompt), projectID, region, modelName)
	if err != nil {
		log.Fatal(err)
	}

}

// countTokens returns the number of tokens for this prompt.
func countTokens(w io.Writer, prompt, projectID, location, modelName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	resp, err := model.CountTokens(ctx, genai.Text(prompt))
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Number of tokens for the prompt: %s\n", strconv.FormatInt(int64(resp.TotalTokens), 10))

	return nil
}
