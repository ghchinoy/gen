package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ghchinoy/gen/internal/model"
)

func init() {
	rootCmd.AddCommand(modelsCmd)
}

var modelsCmd = &cobra.Command{
	Use:     "models",
	Aliases: []string{"m"},
	Short:   "list available models",
	Long:    `Lists available models, foundation, tuned, or Model Garden hosted.`,
	Run:     listModels,
	Hidden:  true,
}

func listModels(cmd *cobra.Command, args []string) {

	mdls := model.List()

	for _, v := range mdls {

		fmt.Println(v)
	}

}
