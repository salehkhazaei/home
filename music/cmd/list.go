package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists input/output devices",
	Long: `Using this command you can see what input/output devices you have along with their ids.'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Printf("invalid argument: %s\n", strings.Join(args, " "))
			cmd.Usage()
			return
		}

		ListOutputDevices()
		ListInputDevices()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
