package cmd

import (
	"fmt"

	"github.com/audiobridge/audiobridge/audio"
	"github.com/spf13/cobra"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List available audio input and output devices",
	Long: `Lists all available audio input (microphone) and output (speakers) 
devices on the system. Shows device index, name, and whether it supports 
input, output, or both.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDevices()
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}

func runDevices() {
	fmt.Println("AudioBridge - Available Audio Devices")
	fmt.Println("======================================")
	fmt.Println()

	// List input devices
	fmt.Println("Input Devices (Microphone):")
	fmt.Println("---------------------------")
	inputDevices, err := audio.ListInputDevices()
	if err != nil {
		fmt.Printf("Error listing input devices: %v\n", err)
		return
	}

	if len(inputDevices) == 0 {
		fmt.Println("  No input devices found")
	} else {
		for _, dev := range inputDevices {
			fmt.Printf("  [%d] %s\n", dev.Index, dev.Name)
		}
	}
	fmt.Println()

	// List output devices
	fmt.Println("Output Devices (Speakers):")
	fmt.Println("---------------------------")
	outputDevices, err := audio.ListOutputDevices()
	if err != nil {
		fmt.Printf("Error listing output devices: %v\n", err)
		return
	}

	if len(outputDevices) == 0 {
		fmt.Println("  No output devices found")
	} else {
		for _, dev := range outputDevices {
			fmt.Printf("  [%d] %s\n", dev.Index, dev.Name)
		}
	}
}
