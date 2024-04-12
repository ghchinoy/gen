package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ghchinoy/gen/internal/model"
)

//go:embed models
var modellist string

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

// TODO - Need to refactor this to use the tablewriter package
func listModels(cmd *cobra.Command, args []string) {
	// read modellist string via CSV reader
	// r := csv.NewReader(strings.NewReader(modellist))
	// modelarray, err := r.ReadAll()
	// if err != nil {
	// 	log.Printf("unable to read model list: %v", err)
	// 	os.Exit(1)
	// }

	mdls := model.List()

	// fmt.Println("Family:\t	Type:\t	Name:")

	// data := [][]string{}
	for _, v := range mdls {

		fmt.Println(v)

		// if strings.HasPrefix(v[0], "#") {
		// 	continue
		// }
		// data = append(data, []string{
		// 	v[0], // group
		// 	v[1], // type
		// 	v[2], // name
		// })
	}
	// table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Family", "Type", "Name"})
	// table.SetBorder(false)
	// table.SetAlignment(tablewriter.ALIGN_LEFT)
	// table.AppendBulk(data)
	// table.Render()

}
