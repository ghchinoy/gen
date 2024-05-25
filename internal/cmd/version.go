package cmd

import (
	"fmt"

	_ "embed"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

//go:embed version
var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gen",
	Long:  `All software have versions. This is gen's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gen %s\n", version)
	},
}
