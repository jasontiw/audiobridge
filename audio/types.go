package audio

// Sample rates
const (
	SampleRate44100 = 44100
	SampleRate48000 = 48000
)

// Channel modes
const (
	ChannelsMono   = 1
	ChannelsStereo = 2
)

// Device represents an audio device
type Device struct {
	Index             int
	Name              string
	HostAPI           string
	MaxInputChannels  int
	MaxOutputChannels int
	DefaultSampleRate float64
}

// Config represents audio configuration
type Config struct {
	SampleRate int
	Channels   int
	DeviceID   int // -1 for default device
}
