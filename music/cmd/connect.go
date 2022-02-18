package cmd

import (
	"fmt"
	"home/music/devices"
	"home/music/network"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connects one input into another output",
	Long:  `Connects one input into another output.`,
	Run: func(cmd *cobra.Command, args []string) {
		var inputDev devices.IDevice
		var err error

		if strings.Contains(input, ":") {
			inputDev, err = network.Bind(input)
			if err != nil {
				panic(err)
			}
		} else {
			inputDev, err = devices.Instance().FindDeviceById(input)
			if err != nil {
				panic(err)
			}
		}

		if inputDev == nil {
			fmt.Println("input device not found")
			os.Exit(1)
			return
		}

		var outputDev devices.IDevice
		if strings.Contains(output, ":") {
			outputDev, err = network.Connect(output)
			if err != nil {
				panic(err)
			}
		} else {
			outputDev, err = devices.Instance().FindDeviceById(output)
			if err != nil {
				panic(err)
			}
		}

		if outputDev == nil {
			fmt.Println("output device not found")
			os.Exit(1)
			return
		}

		samplingChan := make(chan []byte, 1)
		doneChan := make(chan interface{})

		go func() {
			err = inputDev.Capture(sampleRate, channels, samplingChan, doneChan)
			if err != nil {
				panic(err)
			}
		}()

		go func() {
			err = outputDev.Play(channels, sampleRate, samplingChan, doneChan)
			if err != nil {
				fmt.Println("Error on playing", err)
				panic(err)
			}
		}()

		fmt.Println("Press Enter to quit...")
		fmt.Scanln()
		doneChan <- "DONE"
		time.Sleep(time.Second)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringVarP(&input, "input", "i", "", "the input device id or network address (in format of host:port)")
	connectCmd.Flags().StringVarP(&output, "output", "o", "", "the output device id or network address (in format of host:port)")
	connectCmd.Flags().Uint32VarP(&channels, "channels", "c", 1, "no of channels on record/play")
	connectCmd.Flags().Uint32VarP(&sampleRate, "sampleRate", "s", 44800, "sample rate of record/play")
}
