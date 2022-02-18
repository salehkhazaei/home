package cmd

import (
	"fmt"
	"home/music/devices"

	"github.com/spf13/cobra"
)

// outputsCmd represents the outputs command
var outputsCmd = &cobra.Command{
	Use:   "outputs",
	Short: "lists all output devices",
	Long: `lists all output devices`,
	Run: func(cmd *cobra.Command, args []string) {
		ListOutputDevices()
	},
}

func init() {
	listCmd.AddCommand(outputsCmd)
}

func ListOutputDevices() {
	devs, err := devices.Instance().ListPlaybackDevices()
	if err != nil {
		panic(err)
	}

	for _, dev := range devs {
		fmt.Println(dev)
	}
}