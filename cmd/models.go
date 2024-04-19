package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
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
	models, err := model.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	if Outputtype == "json" {
		jsonBytes, err := json.Marshal(models)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(jsonBytes))
	} else {
		data := [][]string{}
		for _, v := range models {
			data = append(data, []string{
				v.Family,
				v.Mode,
				v.Name,
			})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Family", "Mode", "Model ID"})
		table.SetBorder(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.AppendBulk(data)
		table.Render()
	}
}
