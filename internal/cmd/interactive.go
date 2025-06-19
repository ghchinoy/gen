package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ghchinoy/gen/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(interactiveCmd)

	interactiveCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "gemini-2.5-flash", "model name")
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   "Interactive mode",
	Long:    `Interactive mode is a chat mode where you can interact with the model.`,
	RunE:    interactiveMode,
}

func interactiveMode(cmd *cobra.Command, args []string) error {
	fmt.Println("entering interactive mode")
	fmt.Println("type 'exit' or 'quit' to exit")
	fmt.Printf("model: %s\n", modelName)

	cfg := model.Config{
		ProjectID:  projectID,
		RegionID:   region,
		ConfigFile: cfgFile,
		OutputType: Outputtype,
		LogType:    Logtype,
	}

	ctx := context.Background()

	client, err := model.NewClient(ctx, cfg, modelName)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	for {
		fmt.Print("? ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		var buf bytes.Buffer

		// quit | exit
		if strings.EqualFold(input.Text(), "quit") || strings.EqualFold(input.Text(), "exit") {
			return nil
		}

		err := client.GenerateContent(ctx, &buf, input.Text(), nil)
		if err != nil {
			fmt.Printf("error generating content: %v\n", err)
		}

		fmt.Printf("%s\n\n", buf.String())
	}
}
