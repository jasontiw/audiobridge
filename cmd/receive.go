package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/audiobridge/audiobridge/audio"
	"github.com/audiobridge/audiobridge/codec"
	"github.com/audiobridge/audiobridge/jitter"
	"github.com/audiobridge/audiobridge/transport"
	"github.com/spf13/cobra"
)

var (
	receivePort         int
	receiveJitterBuffer int
	receiveDevice       int
)

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive and play audio from the network",
	Long: `Receives audio over UDP from a remote computer and plays it through 
the local audio output device.

Examples:
  audiobridge receive
  audiobridge receive --port 5000 --jitter-buffer 80`,
	RunE: runReceive,
}

func init() {
	// Add receive command to root
	rootCmd.AddCommand(receiveCmd)

	// Define flags
	receiveCmd.Flags().IntVarP(&receivePort, "port", "p", 9876,
		"Listen UDP port")
	receiveCmd.Flags().IntVarP(&receiveJitterBuffer, "jitter-buffer", "j", 50,
		"Jitter buffer latency in milliseconds (default: 50)")
	receiveCmd.Flags().IntVarP(&receiveDevice, "device", "d", -1,
		"Output device index (-1 for default)")
}

func runReceive(cmd *cobra.Command, args []string) error {
	// Validate jitter buffer
	if receiveJitterBuffer < 10 || receiveJitterBuffer > 500 {
		return fmt.Errorf("invalid jitter buffer: %d (must be between 10 and 500 ms)", receiveJitterBuffer)
	}

	fmt.Printf("Starting audio reception...\n")
	fmt.Printf("  Listen Port: %d\n", receivePort)
	fmt.Printf("  Jitter Buffer: %d ms\n", receiveJitterBuffer)
	fmt.Printf("  Output Device: %d\n", receiveDevice)
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

	// Create UDP receiver
	receiver, err := transport.NewUDPReceiver(receivePort)
	if err != nil {
		return fmt.Errorf("failed to create UDP receiver: %w", err)
	}
	defer receiver.Close()
	fmt.Printf("Listening on port %d\n", receivePort)

	// Create Opus decoder
	// Note: We assume 48kHz mono for now, but in production should detect from stream
	dec, err := codec.NewDecoder(48000, 1)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}
	defer dec.Close()
	fmt.Println("Opus decoder initialized")

	// Create jitter buffer
	jbuf, err := jitter.NewBuffer(receiveJitterBuffer)
	if err != nil {
		return fmt.Errorf("failed to create jitter buffer: %w", err)
	}
	defer jbuf.Close()
	fmt.Printf("Jitter buffer initialized (%d ms)\n", receiveJitterBuffer)

	// Create audio player
	player := audio.NewPortAudioPlayer()
	if err := player.Start(receiveDevice, 48000, 1); err != nil {
		return fmt.Errorf("failed to start audio player: %w", err)
	}
	defer player.Close()
	fmt.Println("Audio player started")

	// Start the receive → decode → jitter → playback pipeline
	fmt.Println("\nReceiving audio... Press Ctrl+C to stop.")

	// Start receiver goroutine
	receiveErrChan := make(chan error, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Receive packet
				packet, _, err := receiver.Receive()
				if err != nil {
					receiveErrChan <- err
					continue
				}

				// Push to jitter buffer
				if err := jbuf.Push(packet.Seq, packet.Payload); err != nil {
					fmt.Printf("Error pushing to buffer: %v\n", err)
				}
			}
		}
	}()

	// Main loop: pop from jitter buffer, decode, play
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped.")
			return nil
		case err := <-receiveErrChan:
			fmt.Printf("Receive error: %v\n", err)
		default:
			// Pop from jitter buffer
			opusData, err := jbuf.Pop()
			if err != nil {
				fmt.Printf("Buffer error: %v\n", err)
				continue
			}

			if opusData == nil {
				// Silent frame (packet loss)
				continue
			}

			// Decode
			pcmData, err := dec.Decode(opusData)
			if err != nil {
				fmt.Printf("Decode error: %v\n", err)
				continue
			}

			// Play
			if err := player.Write(pcmData); err != nil {
				fmt.Printf("Playback error: %v\n", err)
				continue
			}
		}
	}
}
