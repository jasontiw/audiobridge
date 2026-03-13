//go:build !cgo

package audio

import (
	"errors"
	"fmt"
)

// ErrNoPortAudio is returned when PortAudio is not available
var ErrNoPortAudio = errors.New("PortAudio not available: build with CGO_ENABLED=1 or install PortAudio")

// ListInputDevices returns all available audio input devices (stub for no_cgo)
func ListInputDevices() ([]Device, error) {
	return nil, ErrNoPortAudio
}

// ListOutputDevices returns all available audio output devices (stub for no_cgo)
func ListOutputDevices() ([]Device, error) {
	return nil, ErrNoPortAudio
}

// PortAudioCapture is a stub when CGO is not available
type PortAudioCapture struct {
	started bool
}

// NewPortAudioCapture creates a new PortAudio capture instance (stub)
func NewPortAudioCapture() *PortAudioCapture {
	return &PortAudioCapture{}
}

// Start begins capture (stub - always returns error)
func (c *PortAudioCapture) Start(deviceID int, sampleRate int, channels int) error {
	return fmt.Errorf("PortAudio not available: %w", ErrNoPortAudio)
}

// Read returns the next audio frame (stub)
func (c *PortAudioCapture) Read() ([]float32, error) {
	return nil, fmt.Errorf("PortAudio not available: %w", ErrNoPortAudio)
}

// Close stops capture (stub)
func (c *PortAudioCapture) Close() error {
	c.started = false
	return nil
}

// PortAudioPlayer is a stub when CGO is not available
type PortAudioPlayer struct {
	started bool
}

// NewPortAudioPlayer creates a new PortAudio player instance (stub)
func NewPortAudioPlayer() *PortAudioPlayer {
	return &PortAudioPlayer{}
}

// Start begins playback (stub - always returns error)
func (p *PortAudioPlayer) Start(deviceID int, sampleRate int, channels int) error {
	return fmt.Errorf("PortAudio not available: %w", ErrNoPortAudio)
}

// Write plays audio (stub)
func (p *PortAudioPlayer) Write(frame []float32) error {
	return fmt.Errorf("PortAudio not available: %w", ErrNoPortAudio)
}

// Close stops playback (stub)
func (p *PortAudioPlayer) Close() error {
	p.started = false
	return nil
}
