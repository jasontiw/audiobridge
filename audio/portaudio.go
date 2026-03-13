//go:build cgo

package audio

import (
	"fmt"
	"sync"

	"github.com/gordonklaus/portaudio"
)

// PortAudioCapture implements AudioCapture using PortAudio
type PortAudioCapture struct {
	stream       *portaudio.Stream
	buffer       []float32
	closed       bool
	done         chan struct{}
	mu           sync.Mutex
	callbackData []float32
}

// NewPortAudioCapture creates a new PortAudio capture instance
func NewPortAudioCapture() *PortAudioCapture {
	return &PortAudioCapture{
		done: make(chan struct{}),
	}
}

// Start begins capture from the specified device
func (c *PortAudioCapture) Start(deviceID int, sampleRate int, channels int) error {
	if c.stream != nil {
		return fmt.Errorf("capture already started")
	}

	var inputDevice *portaudio.DeviceInfo
	var err error

	if deviceID < 0 {
		// Use default input device
		inputDevice, err = portaudio.DefaultInputDevice()
		if err != nil {
			return fmt.Errorf("failed to get default input device: %w", err)
		}
	} else {
		devices, err := portaudio.Devices()
		if err != nil {
			return fmt.Errorf("failed to get devices: %w", err)
		}
		if deviceID >= len(devices) {
			return fmt.Errorf("device index out of range: %d", deviceID)
		}
		inputDevice = devices[deviceID]
	}

	// Calculate buffer size for ~10ms of audio
	framesPerBuffer := sampleRate / 100 // 10ms
	c.callbackData = make([]float32, framesPerBuffer*channels)

	// Use the callback-based API
	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   inputDevice,
			Channels: channels,
			Latency:  inputDevice.DefaultLowInputLatency,
		},
		SampleRate:      float64(sampleRate),
		FramesPerBuffer: framesPerBuffer,
	}, func(in, out []float32) {
		// Copy input data to callback buffer
		c.mu.Lock()
		if len(in) > len(c.callbackData) {
			in = in[:len(c.callbackData)]
		}
		copy(c.callbackData, in)
		c.mu.Unlock()
	})
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}

	c.stream = stream

	if err := stream.Start(); err != nil {
		stream.Close()
		c.stream = nil
		return fmt.Errorf("failed to start audio stream: %w", err)
	}

	return nil
}

// Read returns the next audio frame
func (c *PortAudioCapture) Read() ([]float32, error) {
	if c.stream == nil || c.closed {
		return nil, fmt.Errorf("capture not started")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Copy data from callback buffer
	result := make([]float32, len(c.callbackData))
	copy(result, c.callbackData)

	return result, nil
}

// Close stops capture and releases resources
func (c *PortAudioCapture) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stream == nil {
		return nil
	}

	c.closed = true
	close(c.done)

	err := c.stream.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop audio stream: %w", err)
	}

	err = c.stream.Close()
	c.stream = nil

	if err != nil {
		return fmt.Errorf("failed to close audio stream: %w", err)
	}

	return nil
}
