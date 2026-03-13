//go:build cgo

package audio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// PortAudioPlayer implements AudioPlayer using PortAudio
type PortAudioPlayer struct {
	stream *portaudio.Stream
	closed bool
	buffer []float32
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

	var outputDevice *portaudio.DeviceInfo
	var err error

	if deviceID < 0 {
		// Use default output device
		outputDevice, err = portaudio.DefaultOutputDevice()
		if err != nil {
			return fmt.Errorf("failed to get default output device: %w", err)
		}
	} else {
		devices, err := portaudio.Devices()
		if err != nil {
			return fmt.Errorf("failed to get devices: %w", err)
		}
		if deviceID >= len(devices) {
			return fmt.Errorf("device index out of range: %d", deviceID)
		}
		outputDevice = devices[deviceID]
	}

	// Calculate buffer size for ~10ms of audio
	framesPerBuffer := sampleRate / 100 // 10ms
	p.buffer = make([]float32, framesPerBuffer*channels)

	stream, err := portaudio.OpenStream(p.buffer, portaudio.StreamParameters{
		Output: portaudio.StreamDeviceParameters{
			Device:   outputDevice,
			Channels: channels,
			Latency:  outputDevice.DefaultLowOutputLatency,
		},
		SampleRate:      float64(sampleRate),
		FramesPerBuffer: framesPerBuffer,
	})
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}

	p.stream = stream

	if err := stream.Start(); err != nil {
		stream.Close()
		p.stream = nil
		return fmt.Errorf("failed to start audio stream: %w", err)
	}

	return nil
}

// Write plays the provided audio frame
func (p *PortAudioPlayer) Write(frame []float32) error {
	if p.stream == nil || p.closed {
		return fmt.Errorf("playback not started")
	}

	// Copy frame data to buffer (will be written on next callback)
	if len(frame) > len(p.buffer) {
		frame = frame[:len(p.buffer)]
	}
	copy(p.buffer, frame)

	return nil
}

// Close stops playback and releases resources
func (p *PortAudioPlayer) Close() error {
	if p.stream == nil {
		return nil
	}

	p.closed = true
	err := p.stream.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop audio stream: %w", err)
	}

	err = p.stream.Close()
	p.stream = nil

	if err != nil {
		return fmt.Errorf("failed to close audio stream: %w", err)
	}

	return nil
}
