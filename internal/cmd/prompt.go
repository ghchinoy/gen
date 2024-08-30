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

	promptCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-1.5-flash", "model name")
	promptCmd.PersistentFlags().StringVarP(&modelConfigFile, "config", "c", "", "model parameters")
	promptCmd.PersistentFlags().StringVarP(&promptFile, "file", "f", "", "prompt from file")
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

	var prompt []string
	if promptFile != "" {
		promptBytes, err := os.ReadFile(promptFile)
		if err != nil {
			log.Fatalf("unable to read file %s", promptFile)
		}
		prompt = append(prompt, string(promptBytes))
	} else {
		if len(args) == 0 {
			fmt.Println("please provide prompt")
			log.Fatal("please provide prompt")
		}
		prompt = args
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
		log.Printf("prompt: %s", prompt)
	}

	// Lookup the model based on the name
	m, err := model.Get(modelName)
	if err != nil {
		log.Fatalf("model '%s' is not supported\n", modelName)
	}

	ctx := context.Background()

	err = m.Use(ctx, cfg, prompt)
	if err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}

}
