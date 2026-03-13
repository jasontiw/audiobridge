//go:build !cgo

package audio

import (
	"errors"
)

// ErrNoPortAudio is returned when PortAudio is not available
var ErrNoPortAudio = errors.New("PortAudio not available: build with CGO_ENABLED=1 or install PortAudio")

// ListInputDevices returns all available audio input devices (stub for no_portaudio)
func ListInputDevices() ([]Device, error) {
	return nil, ErrNoPortAudio
}

// ListOutputDevices returns all available audio output devices (stub for no_portaudio)
func ListOutputDevices() ([]Device, error) {
	return nil, ErrNoPortAudio
}
