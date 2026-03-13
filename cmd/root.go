package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "audiobridge",
	Short: "Share audio between computers on the same network",
	Long: `AudioBridge is a cross-platform tool for sharing audio between 
computers on the same local network. Works on Windows, macOS, and Linux.

Usage:
  audiobridge send --target 192.168.1.10   # Send audio to another PC
  audiobridge receive                       # Receive and play audio
  audiobridge devices                       # List available audio devices
  audiobridge status                        # Show connection status

For more information, visit: https://github.com/jasontiw/audiobridge`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info",
		"Log level: debug, info, warn, error")
}

func setupLogging() {
	// Validate log level
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	level := strings.ToLower(logLevel)
	if !validLevels[level] {
		fmt.Fprintf(os.Stderr, "Invalid log level: %s. Valid levels: debug, info, warn, error\n", logLevel)
		os.Exit(1)
	}

	// For now, we use the standard log package
	// Future: integrate with structured logging
	_ = level // Will be used when logging is implemented
}
