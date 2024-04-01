package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile    string
	region     string
	projectID  string
	outputtype string
	logtype    string

	rootCmd = &cobra.Command{
		Use:   "gen",
		Short: "access generative ai on google cloud",
		Long:  `gen is a command-line tool for Google Cloud hosted generative ai models - foundation, tuned, and Model Garden models.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/gen/gen.yaml)")
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "Google Cloud Project ID")
	//rootCmd.MarkPersistentFlagRequired("project")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "region for generative AI endpoint")
	rootCmd.PersistentFlags().StringVar(&outputtype, "output", "text", "output type")

	rootCmd.PersistentFlags().StringVar(&logtype, "log", "none", "logging output")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name "cxctl" (without extension).
		viper.SetConfigName("gen") // name of config file (without extension)
		viper.AddConfigPath(os.Getenv("HOME") + "/.config/gen")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("gen")
	}

	//viper.SetEnvPrefix("GEN")
	viper.AutomaticEnv()

	if viper.IsSet("PROJECT_ID") {
		rootCmd.Flags().Set("project", fmt.Sprintf("%v", viper.Get("PROJECT_ID")))
	}
	if viper.IsSet("REGION") {
		rootCmd.Flags().Set("region", fmt.Sprintf("%v", viper.Get("REGION")))
	}

	// check if there are prefixed env vars and bind them
	// ref. h/t to https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
