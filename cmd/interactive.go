package cmd

import (
	"log"

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
	log.Printf("interactive mode tbd")
}
