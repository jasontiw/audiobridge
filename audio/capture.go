package audio

// AudioCapture defines the interface for capturing audio from input devices
type AudioCapture interface {
	// Start begins capture from the specified device
	// deviceID: -1 for default, or specific device index
	// sampleRate: 44100 or 48000
	// channels: 1 (mono) or 2 (stereo)
	Start(deviceID int, sampleRate int, channels int) error

	// Read returns the next audio frame
	// Returns nil when capture is closed
	Read() ([]float32, error)

	// Close stops capture and releases resources
	Close() error
}

// CaptureSource defines the type of audio source
type CaptureSource int

const (
	CaptureSourceMicrophone CaptureSource = iota
	CaptureSourceSystem
)

// CaptureConfig holds configuration for audio capture
type CaptureConfig struct {
	Source     CaptureSource
	DeviceID   int
	SampleRate int
	Channels   int
}
