
package cmd

import (
	"fmt"
	"home/music/devices"

	"github.com/spf13/cobra"
)

// inputsCmd represents the inputs command
var inputsCmd = &cobra.Command{
	Use:   "inputs",
	Short: "lists all input devices",
	Long: `lists all input devices`,
	Run: func(cmd *cobra.Command, args []string) {
		ListInputDevices()
	},
}

func init() {
	listCmd.AddCommand(inputsCmd)
}

func ListInputDevices() {
	devs, err := devices.Instance().ListCaptureDevices()
	if err != nil {
		panic(err)
	}

	for _, dev := range devs {
		fmt.Println(dev)
	}
}