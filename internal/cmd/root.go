package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	rootCmd.MarkPersistentFlagRequired("project")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "region for generative AI endpoint")
	rootCmd.MarkPersistentFlagRequired("region")
	rootCmd.PersistentFlags().StringVar(&Outputtype, "output", "text", "output type")
	rootCmd.PersistentFlags().StringVar(&Logtype, "log", "none", "logging output")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name "gen" (without extension).
		viper.SetConfigName("gen") // name of config file (without extension)
		viper.AddConfigPath(os.Getenv("HOME") + "/.config/gen")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
	}

	// bind environment variables
	//viper.SetEnvPrefix("GEN") // env variables prefix
	viper.AutomaticEnv()

	if viper.IsSet("PROJECT_ID") {
		project_flag, _ := rootCmd.Flags().GetString("project")
		if project_flag == "" {
			rootCmd.Flags().Set("project", fmt.Sprintf("%v", viper.Get("PROJECT_ID")))
		}
	}
	if viper.IsSet("REGION") {
		region_flag, _ := rootCmd.Flags().GetString("region")
		if region_flag == "" {
			rootCmd.Flags().Set("region", fmt.Sprintf("%v", viper.Get("REGION")))
		}
	}

	// check if there are prefixed env vars and bind them
	// ref. h/t to https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
