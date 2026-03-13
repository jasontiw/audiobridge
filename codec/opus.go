package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Opus encoder/decoder implementation
// Note: This is a placeholder implementation that passes audio through without compression.
// For production, configure hraban/opus with proper CGO/LibOpus build environment.

var (
	// ErrNotInitialized is returned when encoder/decoder is not initialized
	ErrNotInitialized = errors.New("encoder/decoder not initialized")

	// ErrOpusNotAvailable is returned when opus is not properly configured
	ErrOpusNotAvailable = errors.New("Opus codec not available - configure CGO with libopus")
)

// EncoderConfig holds configuration for Opus encoder
type EncoderConfig struct {
	SampleRate  int    // 48000 or 44100
	Channels    int    // 1 = mono, 2 = stereo
	Bitrate     int    // bits per second (e.g., 64000 for 64 kbps)
	Application string // "voip", "audio", "lowdelay"
}

// DefaultEncoderConfig returns a default encoder configuration for voice
func DefaultEncoderConfig() EncoderConfig {
	return EncoderConfig{
		SampleRate:  48000,
		Channels:    1,
		Bitrate:     64000,
		Application: "voip",
	}
}

// OpusEncoder wraps the opus encoder
type OpusEncoder struct {
	encoder interface{} // *opus.Encoder - nil when using stub
	config  EncoderConfig
	// For the stub: we pass audio through but encode as a simple format
	frameCount uint32
}

// NewEncoder creates a new Opus encoder
func NewEncoder(config EncoderConfig) (*OpusEncoder, error) {
	// Try to use hraban/opus if available, otherwise use stub
	// The actual opus library requires CGO with libopus
	// For now, we create a pass-through encoder

	return &OpusEncoder{
		config: config,
	}, nil
}

// Encode encodes PCM audio (float32) to Opus format
// Note: This is a stub that wraps PCM in a simple container format
// Real implementation would use opus.Encoder
func (e *OpusEncoder) Encode(pcmData []float32) ([]byte, error) {
	if e.encoder == nil {
		// Stub: pass through with simple framing
		// This allows the application to work without opus libraries

		// Convert float32 to bytes
		pcmBytes := make([]byte, len(pcmData)*4)
		for i, sample := range pcmData {
			// Clamp to [-1, 1]
			if sample > 1 {
				sample = 1
			} else if sample < -1 {
				sample = -1
			}
			bits := int32(sample * float32(math.MaxInt16))
			binary.LittleEndian.PutUint32(pcmBytes[i*4:], uint32(bits))
		}

		// Simple packet format:
		// [4 bytes: sequence][4 bytes: sample rate][2 bytes: channels][2 bytes: frame size][PCM data]
		header := make([]byte, 12)
		binary.LittleEndian.PutUint32(header[0:4], e.frameCount)
		binary.LittleEndian.PutUint32(header[4:8], uint32(e.config.SampleRate))
		binary.LittleEndian.PutUint16(header[8:10], uint16(e.config.Channels))
		binary.LittleEndian.PutUint16(header[10:12], uint16(len(pcmData)))

		e.frameCount++

		// Combine header and data
		result := make([]byte, 12+len(pcmBytes))
		copy(result[0:12], header)
		copy(result[12:], pcmBytes)

		return result, nil
	}

	return nil, ErrOpusNotAvailable
}

// SetBitrate changes the encoder bitrate
func (e *OpusEncoder) SetBitrate(bitrate int) error {
	e.config.Bitrate = bitrate
	return nil
}

// Close closes the encoder and releases resources
func (e *OpusEncoder) Close() error {
	e.encoder = nil
	return nil
}

// OpusDecoder wraps the opus decoder
type OpusDecoder struct {
	decoder    interface{} // *opus.Decoder - nil when using stub
	sampleRate int
	channels   int
	frameCount uint32
}

// NewDecoder creates a new Opus decoder
func NewDecoder(sampleRate, channels int) (*OpusDecoder, error) {
	return &OpusDecoder{
		sampleRate: sampleRate,
		channels:   channels,
	}, nil
}

// Decode decodes Opus audio to PCM format (float32)
// Note: This is a stub that unwraps the simple container format
func (d *OpusDecoder) Decode(opusData []byte) ([]float32, error) {
	if d.decoder == nil {
		// Stub: unpack the simple format
		if len(opusData) < 12 {
			return nil, fmt.Errorf("invalid packet: too short")
		}

		// Read header
		seq := binary.LittleEndian.Uint32(opusData[0:4])
		sampleRate := binary.LittleEndian.Uint32(opusData[4:8])
		channels := binary.LittleEndian.Uint16(opusData[8:10])
		frameSize := binary.LittleEndian.Uint16(opusData[10:12])

		// Validate
		if sampleRate != uint32(d.sampleRate) || channels != uint16(d.channels) {
			return nil, fmt.Errorf("format mismatch: expected %dHz %dch, got %dHz %dch",
				d.sampleRate, d.channels, sampleRate, channels)
		}

		// Skip header and read PCM
		pcmBytes := opusData[12:]
		expectedBytes := int(frameSize) * 4
		if len(pcmBytes) < expectedBytes {
			return nil, fmt.Errorf("invalid packet: expected %d bytes, got %d", expectedBytes, len(pcmBytes))
		}

		// Convert bytes to float32
		pcmData := make([]float32, frameSize)
		for i := 0; i < int(frameSize); i++ {
			bits := binary.LittleEndian.Uint32(pcmBytes[i*4 : i*4+4])
			pcmData[i] = float32(int32(bits)) / float32(math.MaxInt16)
		}

		d.frameCount = seq
		return pcmData, nil
	}

	return nil, ErrOpusNotAvailable
}

// Close closes the decoder and releases resources
func (d *OpusDecoder) Close() error {
	d.decoder = nil
	return nil
}

// Encode encodes PCM audio to Opus format (convenience function)
// Uses default encoder configuration
func Encode(pcmData []float32) ([]byte, error) {
	enc, err := NewEncoder(DefaultEncoderConfig())
	if err != nil {
		return nil, err
	}
	defer enc.Close()

	return enc.Encode(pcmData)
}

// Decode decodes Opus audio to PCM format (convenience function)
func Decode(opusData []byte) ([]float32, error) {
	dec, err := NewDecoder(48000, 1)
	if err != nil {
		return nil, err
	}
	defer dec.Close()

	return dec.Decode(opusData)
}
