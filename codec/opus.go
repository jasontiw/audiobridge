package codec

import (
	"errors"
)

// Opus encoder/decoder stub
// Full implementation will be added in Phase 2

var (
	// ErrNotImplemented is returned when Opus encoding/decoding is called
	ErrNotImplemented = errors.New("Opus codec not implemented yet - coming in Phase 2")
)

// Encode encodes PCM audio to Opus format
// This is a stub - full implementation in Phase 2
func Encode(pcmData []float32) ([]byte, error) {
	// TODO: Implement Opus encoding in Phase 2
	return nil, ErrNotImplemented
}

// Decode decodes Opus audio to PCM format
// This is a stub - full implementation in Phase 2
func Decode(opusData []byte) ([]float32, error) {
	// TODO: Implement Opus decoding in Phase 2
	return nil, ErrNotImplemented
}

// EncoderConfig holds configuration for Opus encoder
type EncoderConfig struct {
	SampleRate int
	Channels   int
	Bitrate    int // bits per second
}

// NewEncoder creates a new Opus encoder (stub)
func NewEncoder(config EncoderConfig) (*OpusEncoder, error) {
	// TODO: Implement in Phase 2
	return nil, ErrNotImplemented
}

// OpusEncoder is a stub for the Opus encoder
type OpusEncoder struct {
	config EncoderConfig
}

// Encode encodes audio (stub)
func (e *OpusEncoder) Encode(pcmData []float32) ([]byte, error) {
	return nil, ErrNotImplemented
}

// Close closes the encoder (stub)
func (e *OpusEncoder) Close() error {
	return nil
}

// NewDecoder creates a new Opus decoder (stub)
func NewDecoder(sampleRate, channels int) (*OpusDecoder, error) {
	// TODO: Implement in Phase 2
	return nil, ErrNotImplemented
}

// OpusDecoder is a stub for the Opus decoder
type OpusDecoder struct {
	sampleRate int
	channels   int
}

// Decode decodes audio (stub)
func (d *OpusDecoder) Decode(opusData []byte) ([]float32, error) {
	return nil, ErrNotImplemented
}

// Close closes the decoder (stub)
func (d *OpusDecoder) Close() error {
	return nil
}
