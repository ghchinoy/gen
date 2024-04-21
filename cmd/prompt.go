package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ghchinoy/gen/internal/model"
	"github.com/spf13/cobra"
)

var (
	systemInstructions string
)

func init() {
	rootCmd.AddCommand(promptCmd)

	promptCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-1.0-pro", "model name")
	promptCmd.PersistentFlags().StringVarP(&modelConfigFile, "config", "c", "", "model parameters")
}

var promptCmd = &cobra.Command{
	Use:     "prompt",
	Aliases: []string{"p"},
	Short:   "Prompt a model",
	Long:    `Provide prompt parts to a model to generate content`,
	// Run:     generateContentForModel,
	Run: generateContent,
}

// generateContent prompts a model to generate content based on the provided prompt.
func generateContent(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("please provide prompt")
		log.Fatal("please provide prompt")
	}
	//log.Printf("ProjectID: %s, Region: %s\n", projectID, region)

	cfgB := model.ConfigBuilder{}

	// Set the model configuration
	cfgB.ProjectID(projectID).RegionID(region).ConfigFile(cfgFile).OutputType(Outputtype).LogType(Logtype)
	// fmt.Printf("Model config: %s\n", cfgB.Describe())
	cfg, err := cfgB.Build()
	log.Printf("Model Config: %+v", cfg)

	if err != nil {
		log.Fatalf("error building config: %v", err)
	}

	if Logtype != "none" {
		log.Printf("model: %s", modelName)
		log.Printf("prompt: %s", args)
	}

	// Lookup the model based on the name
	m, err := model.Get(modelName)
	if err != nil {
		log.Fatalf("model '%s' is not supported\n", modelName)
	}

	ctx := context.Background()

	err = m.Use(ctx, cfg, args)
	if err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}

}
