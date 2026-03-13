package audio

// AudioPlayer defines the interface for playing audio to output devices
type AudioPlayer interface {
	// Start begins playback to the specified device
	Start(deviceID int, sampleRate int, channels int) error

	// Write plays the provided audio frame
	Write(frame []float32) error

	// Close stops playback and releases resources
	Close() error
}

// PlaybackConfig holds configuration for audio playback
type PlaybackConfig struct {
	DeviceID   int
	SampleRate int
	Channels   int
}
