package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/vertexai/genai"
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
	Short:   "interactive mode",
	Long:    `Interactive mode is a chat mode where you can interact with the model.`,
	Run:     interactiveMode,
}

func interactiveMode(cmd *cobra.Command, args []string) {
	log.Printf("entering interactive mode")
	for {
		fmt.Print("? ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		//fmt.Printf("You asked: %s", input.Text())
		var buf bytes.Buffer
		prompt := genai.Text(input.Text())
		if err := generateContentGemini(&buf, projectID, region, modelName, []genai.Part{prompt}); err != nil {
			log.Printf("error generating content: %v", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n\n", buf.String())
	}
}
