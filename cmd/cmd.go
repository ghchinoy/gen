package cmd

import "github.com/spf13/cobra"

var (
	// TODO - Look at ways to remove the need to export these two variable outside the package
	modelName       string
	modelConfigFile string
	//modelConfig     map[string]interface{}

	// Used for flags.
	cfgFile   string
	region    string
	projectID string
	// TODO - Look for ways to remove the need to export this outside of package
	Outputtype string
	// TODO - Look for ways to remove the need to export this outside of package
	Logtype string

	rootCmd = &cobra.Command{
		Use:   "gen",
		Short: "access generative ai on google cloud",
		Long:  `gen is a command-line tool for Google Cloud hosted generative ai models - foundation, tuned, and Model Garden models.`,
	}
)
