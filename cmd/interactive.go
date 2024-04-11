package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/ghchinoy/gen/internal/model"
	"github.com/spf13/cobra"
)

var (
// modelName string // for reference, defined in prompt
)

func init() {
	rootCmd.AddCommand(interactiveCmd)

	interactiveCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-1.0-pro", "model name")
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   "Interactive mode",
	Long:    `Interactive mode is a chat mode where you can interact with the model.`,
	Run:     interactiveMode,
}

func interactiveMode(cmd *cobra.Command, args []string) {
	log.Print("entering interactive mode")
	log.Print("type 'exit' or 'quit' to exit")
	log.Printf("model: %s", modelName)

	if !strings.Contains(modelName, "gemini") {
		log.Print("Apologies, only gemini models at this time")
		os.Exit(0)
	}
	for {
		fmt.Print("? ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		var buf bytes.Buffer

		// quit | exit
		if strings.EqualFold(input.Text(), "quit") || strings.EqualFold(input.Text(), "exit") {
			os.Exit(0)
		}

		// gemini
		prompt := genai.Text(input.Text())
		if err := model.GenerateContentGemini(&buf, projectID, region, modelName, []genai.Part{prompt}); err != nil {
			log.Printf("error generating content: %v", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n\n", buf.String())
	}
}
