package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ghchinoy/gen/internal/model"
	"github.com/spf13/cobra"
)

var (
	systemInstructions string
)

func init() {
	rootCmd.AddCommand(promptCmd)

	promptCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-2.5-flash", "model name")
promptCmd.Flag("model").DefValue = "gemini-2.5-flash"
	//promptCmd.PersistentFlags().StringArrayVarP(&modelNames, "model", "m", []string{"gemini-1.5-flash"}, "model name(s)")
	promptCmd.PersistentFlags().StringVarP(&modelConfigFile, "config", "c", "", "model parameters")
	promptCmd.PersistentFlags().StringVarP(&promptFile, "file", "f", "", "prompt from file")
}

var promptCmd = &cobra.Command{
	Use:     "prompt",
	Aliases: []string{"p"},
	Short:   "Prompt a model",
	Long:    `Provide prompt parts to a model to generate content`,
	RunE:    generateContentE,
}

// generateContentE prompts a model to generate content based on the provided prompt.
func generateContentE(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("model").Changed {
		modelName = "gemini-2.5-flash"
	}

	var prompt string

	if promptFile != "" {
		promptBytes, err := os.ReadFile(promptFile)
		if err != nil {
			return fmt.Errorf("unable to read file %s: %w", promptFile, err)
		}
		prompt = string(promptBytes)
	} else {
		if len(args) == 0 {
			return fmt.Errorf("please provide prompt")
		}
		prompt = strings.Join(args, " ")
	}

	cfg := model.Config{
		ProjectID:  projectID,
		RegionID:   region,
		ConfigFile: cfgFile,
		OutputType: Outputtype,
		LogType:    Logtype,
	}

	if Logtype != "none" {
		fmt.Printf("model: %s\n", modelName)
		fmt.Printf("prompt: %s\n", prompt)
	}

	ctx := context.Background()

	client, err := model.NewClient(ctx, cfg, modelName)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	return client.GenerateContent(ctx, os.Stdout, prompt, nil)
}

