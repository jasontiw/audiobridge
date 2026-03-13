package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/audiobridge/audiobridge/audio"
	"github.com/audiobridge/audiobridge/codec"
	"github.com/audiobridge/audiobridge/transport"
	"github.com/spf13/cobra"
)

var (
	sendSource     string
	sendTarget     string
	sendPort       int
	sendBitrate    int
	sendSampleRate int
	sendChannels   int
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send audio over the network",
	Long: `Captures audio from a microphone or system audio and streams it 
to a target computer over UDP.

Examples:
  audiobridge send --target 192.168.1.10
  audiobridge send --target 192.168.1.10 --port 5000 --bitrate 128000
  audiobridge send --system --target 192.168.1.10`,
	RunE: runSend,
}

func init() {
	// Add send command to root
	rootCmd.AddCommand(sendCmd)

	// Define flags
	sendCmd.Flags().StringVarP(&sendSource, "source", "s", "microphone",
		"Audio source: microphone or system")
	sendCmd.Flags().StringVarP(&sendTarget, "target", "t", "",
		"Target IP address (required)")
	sendCmd.Flags().IntVarP(&sendPort, "port", "p", 9876,
		"Target UDP port")
	sendCmd.Flags().IntVarP(&sendBitrate, "bitrate", "b", 64000,
		"Opus bitrate in bits per second (default: 64000)")
	sendCmd.Flags().IntVarP(&sendSampleRate, "rate", "r", 48000,
		"Audio sample rate (default: 48000)")
	sendCmd.Flags().IntVarP(&sendChannels, "channels", "c", 1,
		"Number of channels: 1=mono, 2=stereo (default: 1)")

	// Mark target as required
	_ = sendCmd.MarkFlagRequired("target")
}

func runSend(cmd *cobra.Command, args []string) error {
	// Validate source
	if sendSource != "microphone" && sendSource != "system" {
		return fmt.Errorf("invalid source: %s (must be 'microphone' or 'system')", sendSource)
	}

	// Validate bitrate
	if sendBitrate < 6000 || sendBitrate > 512000 {
		return fmt.Errorf("invalid bitrate: %d (must be between 6000 and 512000)", sendBitrate)
	}

	// Validate sample rate
	if sendSampleRate != 44100 && sendSampleRate != 48000 {
		return fmt.Errorf("invalid sample rate: %d (must be 44100 or 48000)", sendSampleRate)
	}

	// Validate channels
	if sendChannels < 1 || sendChannels > 2 {
		return fmt.Errorf("invalid channels: %d (must be 1 or 2)", sendChannels)
	}

	fmt.Printf("Starting audio transmission...\n")
	fmt.Printf("  Target: %s:%d\n", sendTarget, sendPort)
	fmt.Printf("  Source: %s\n", sendSource)
	fmt.Printf("  Bitrate: %d bps\n", sendBitrate)
	fmt.Printf("  Sample Rate: %d Hz\n", sendSampleRate)
	fmt.Printf("  Channels: %d\n", sendChannels)
	fmt.Println()

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Create UDP sender
	address := fmt.Sprintf("%s:%d", sendTarget, sendPort)
	sender, err := transport.NewUDPSender(address)
	if err != nil {
		return fmt.Errorf("failed to create UDP sender: %w", err)
	}
	defer sender.Close()
	fmt.Printf("Connected to %s\n", address)

	// Create Opus encoder
	enc, err := codec.NewEncoder(codec.EncoderConfig{
		SampleRate: sendSampleRate,
		Channels:   sendChannels,
		Bitrate:    sendBitrate,
	})
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}
	defer enc.Close()
	fmt.Println("Opus encoder initialized")

	// Create audio capture
	var capture *audio.PortAudioCapture

	if sendSource == "microphone" {
		// Use PortAudio for microphone capture
		capture = audio.NewPortAudioCapture()
	} else {
		// System audio not yet implemented
		return fmt.Errorf("system audio capture not yet implemented")
	}

	if err := capture.Start(-1, sendSampleRate, sendChannels); err != nil {
		return fmt.Errorf("failed to start capture: %w", err)
	}
	defer capture.Close()
	fmt.Printf("Capturing from %s\n", sendSource)

	// Start the capture → encode → send pipeline
	fmt.Println("\nStreaming audio... Press Ctrl+C to stop.")

	// Audio buffer - 10ms frame
	frameSize := sendSampleRate * sendChannels * 10 / 1000
	audioBuffer := make([]float32, frameSize)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped.")
			return nil
		case <-ticker.C:
			// Read audio frame
			frame, err := capture.Read()
			if err != nil {
				fmt.Printf("Error reading audio: %v\n", err)
				continue
			}
			if frame == nil {
				// Capture closed
				fmt.Println("Audio capture closed")
				return nil
			}

			// Ensure we have the right size
			if len(frame) != len(audioBuffer) {
				frameSize = len(frame)
				audioBuffer = make([]float32, frameSize)
			}
			copy(audioBuffer, frame)

			// Encode
			opusData, err := enc.Encode(audioBuffer)
			if err != nil {
				fmt.Printf("Error encoding: %v\n", err)
				continue
			}

			// Send
			if err := sender.Send(opusData); err != nil {
				fmt.Printf("Error sending: %v\n", err)
				continue
			}
		}
	}
}
