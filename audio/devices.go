//go:build !no_portaudio && cgo

package audio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// ListInputDevices returns all available audio input devices
func ListInputDevices() ([]Device, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize PortAudio: %w", err)
	}
	defer portaudio.Terminate()

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	var inputDevices []Device
	for i, dev := range devices {
		if dev.MaxInputChannels > 0 {
			inputDevices = append(inputDevices, Device{
				Index:             i,
				Name:              dev.Name,
				HostAPI:           dev.HostApi,
				MaxInputChannels:  dev.MaxInputChannels,
				MaxOutputChannels: dev.MaxOutputChannels,
				DefaultSampleRate: dev.DefaultSampleRate,
			})
		}
	}

	return inputDevices, nil
}

// ListOutputDevices returns all available audio output devices
func ListOutputDevices() ([]Device, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize PortAudio: %w", err)
	}
	defer portaudio.Terminate()

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	var outputDevices []Device
	for i, dev := range devices {
		if dev.MaxOutputChannels > 0 {
			outputDevices = append(outputDevices, Device{
				Index:             i,
				Name:              dev.Name,
				HostAPI:           dev.HostApi,
				MaxInputChannels:  dev.MaxInputChannels,
				MaxOutputChannels: dev.MaxOutputChannels,
				DefaultSampleRate: dev.DefaultSampleRate,
			})
		}
	}

	return outputDevices, nil
}
