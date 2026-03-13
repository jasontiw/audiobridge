//go:build cgo

package audio

import (
	"fmt"
	"sync"

	"github.com/gordonklaus/portaudio"
)

// PortAudioPlayer implements AudioPlayer using PortAudio
type PortAudioPlayer struct {
	stream       *portaudio.Stream
	closed       bool
	mu           sync.Mutex
	callbackData []float32
}

// NewPortAudioPlayer creates a new PortAudio player instance
func NewPortAudioPlayer() *PortAudioPlayer {
	return &PortAudioPlayer{}
}

// Start begins playback to the specified device
func (p *PortAudioPlayer) Start(deviceID int, sampleRate int, channels int) error {
	if p.stream != nil {
		return fmt.Errorf("playback already started")
	}

	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %w", err)
	}

	var outputDevice *portaudio.DeviceInfo
	var err error

	if deviceID < 0 {
		// Use default output device
		outputDevice, err = portaudio.DefaultOutputDevice()
		if err != nil {
			portaudio.Terminate()
			return fmt.Errorf("failed to get default output device: %w", err)
		}
	} else {
		devices, err := portaudio.Devices()
		if err != nil {
			portaudio.Terminate()
			return fmt.Errorf("failed to get devices: %w", err)
		}
		if deviceID >= len(devices) {
			portaudio.Terminate()
			return fmt.Errorf("device index out of range: %d", deviceID)
		}
		outputDevice = devices[deviceID]
	}

	// Calculate buffer size for ~10ms of audio
	framesPerBuffer := sampleRate / 100 // 10ms
	p.callbackData = make([]float32, framesPerBuffer*channels)

	// Use the callback-based API
	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Output: portaudio.StreamDeviceParameters{
			Device:   outputDevice,
			Channels: channels,
			Latency:  outputDevice.DefaultLowOutputLatency,
		},
		SampleRate:      float64(sampleRate),
		FramesPerBuffer: framesPerBuffer,
	}, func(in, out []float32) {
		// Copy data from callback buffer to output
		p.mu.Lock()
		copy(out, p.callbackData)
		// Fill with silence for next callback
		for i := range p.callbackData {
			p.callbackData[i] = 0
		}
		p.mu.Unlock()
	})
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}

	p.stream = stream

	if err := stream.Start(); err != nil {
		stream.Close()
		p.stream = nil
		return fmt.Errorf("failed to start playback: %w", err)
	}

	return nil
}

// Write plays the provided audio frame
func (p *PortAudioPlayer) Write(frame []float32) error {
	if p.stream == nil || p.closed {
		return fmt.Errorf("playback not started")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Copy frame data to callback buffer
	if len(frame) > len(p.callbackData) {
		frame = frame[:len(p.callbackData)]
	}
	copy(p.callbackData, frame)

	return nil
}

// Close stops playback and releases resources
func (p *PortAudioPlayer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stream == nil {
		portaudio.Terminate()
		return nil
	}

	p.closed = true
	err := p.stream.Stop()
	if err != nil {
		portaudio.Terminate()
		return fmt.Errorf("failed to stop audio stream: %w", err)
	}

	err = p.stream.Close()
	p.stream = nil
	portaudio.Terminate()

	if err != nil {
		return fmt.Errorf("failed to close audio stream: %w", err)
	}

	return nil
}
