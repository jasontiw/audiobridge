package codec

import (
	"testing"
)

func TestOpusEncoder_NewEncoder(t *testing.T) {
	tests := []struct {
		name    string
		config  EncoderConfig
		wantErr bool
	}{
		{
			name: "default voice config",
			config: EncoderConfig{
				SampleRate:  48000,
				Channels:    1,
				Bitrate:     64000,
				Application: "voip",
			},
			wantErr: false,
		},
		{
			name: "stereo audio",
			config: EncoderConfig{
				SampleRate:  44100,
				Channels:    2,
				Bitrate:     128000,
				Application: "audio",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncoder(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && enc == nil {
				t.Error("NewEncoder() returned nil encoder")
			}
			if enc != nil {
				enc.Close()
			}
		})
	}
}

func TestOpusEncoder_Decoder_EncodeDecode(t *testing.T) {
	// Create encoder
	enc, err := NewEncoder(EncoderConfig{
		SampleRate:  48000,
		Channels:    1,
		Bitrate:     64000,
		Application: "voip",
	})
	if err != nil {
		t.Fatalf("NewEncoder() failed: %v", err)
	}
	defer enc.Close()

	// Create decoder
	dec, err := NewDecoder(48000, 1)
	if err != nil {
		t.Fatalf("NewDecoder() failed: %v", err)
	}
	defer dec.Close()

	// Test data: 480 samples (10ms at 48kHz)
	pcmOriginal := make([]float32, 480)
	for i := range pcmOriginal {
		pcmOriginal[i] = float32(i) / 1000.0 // Simple ramp
	}

	// Encode
	encoded, err := enc.Encode(pcmOriginal)
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatal("Encode() returned empty data")
	}

	// Decode
	pcmDecoded, err := dec.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode() failed: %v", err)
	}

	// Verify
	if len(pcmDecoded) != len(pcmOriginal) {
		t.Errorf("Decoded length: got %d, want %d", len(pcmDecoded), len(pcmOriginal))
	}

	// Check samples are approximately equal (some loss is expected in the format)
	for i := 0; i < len(pcmOriginal) && i < len(pcmDecoded); i++ {
		diff := pcmOriginal[i] - pcmDecoded[i]
		if diff < 0 {
			diff = -diff
		}
		if diff > 0.01 { // Allow 1% tolerance
			t.Errorf("Sample %d: got %v, want ~%v (diff %v)", i, pcmDecoded[i], pcmOriginal[i], diff)
		}
	}
}

func TestOpusEncoder_SetBitrate(t *testing.T) {
	enc, err := NewEncoder(DefaultEncoderConfig())
	if err != nil {
		t.Fatalf("NewEncoder() failed: %v", err)
	}
	defer enc.Close()

	// Change bitrate
	err = enc.SetBitrate(128000)
	if err != nil {
		t.Errorf("SetBitrate() failed: %v", err)
	}

	// Encode should still work
	_, err = enc.Encode(make([]float32, 480))
	if err != nil {
		t.Errorf("Encode() after SetBitrate() failed: %v", err)
	}
}

func TestOpusDecoder_InvalidPacket(t *testing.T) {
	dec, err := NewDecoder(48000, 1)
	if err != nil {
		t.Fatalf("NewDecoder() failed: %v", err)
	}
	defer dec.Close()

	tests := []struct {
		name      string
		packet    []byte
		wantErr   bool
		errString string
	}{
		{
			name:    "empty packet",
			packet:  []byte{},
			wantErr: true,
		},
		{
			name:    "too short packet",
			packet:  []byte{0x00, 0x01, 0x02},
			wantErr: true,
		},
		{
			name:    "format mismatch",
			packet:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // header with wrong format
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dec.Decode(tt.packet)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpusEncoder_Close(t *testing.T) {
	enc, err := NewEncoder(DefaultEncoderConfig())
	if err != nil {
		t.Fatalf("NewEncoder() failed: %v", err)
	}

	// Close should not error
	if err := enc.Close(); err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Encode after close - should still work (stub behavior)
	_, err = enc.Encode(make([]float32, 480))
	if err != nil {
		t.Logf("Note: Encode() after Close() returned error (may be expected): %v", err)
	}
}

func TestDefaultEncoderConfig(t *testing.T) {
	config := DefaultEncoderConfig()

	if config.SampleRate != 48000 {
		t.Errorf("Default sample rate: got %d, want 48000", config.SampleRate)
	}
	if config.Channels != 1 {
		t.Errorf("Default channels: got %d, want 1", config.Channels)
	}
	if config.Bitrate != 64000 {
		t.Errorf("Default bitrate: got %d, want 64000", config.Bitrate)
	}
	if config.Application != "voip" {
		t.Errorf("Default application: got %s, want voip", config.Application)
	}
}

func TestEncodeDecode_ConsecutiveFrames(t *testing.T) {
	// Test encoding/decoding multiple frames
	enc, err := NewEncoder(EncoderConfig{
		SampleRate:  48000,
		Channels:    1,
		Bitrate:     64000,
		Application: "voip",
	})
	if err != nil {
		t.Fatalf("NewEncoder() failed: %v", err)
	}
	defer enc.Close()

	dec, err := NewDecoder(48000, 1)
	if err != nil {
		t.Fatalf("NewDecoder() failed: %v", err)
	}
	defer dec.Close()

	// Encode 10 consecutive frames
	frames := make([][]float32, 10)
	encoded := make([][]byte, 10)

	for i := range frames {
		frames[i] = make([]float32, 480)
		for j := range frames[i] {
			frames[i][j] = float32(i*j) / 1000.0
		}

		var err error
		encoded[i], err = enc.Encode(frames[i])
		if err != nil {
			t.Fatalf("Encode() frame %d failed: %v", i, err)
		}
	}

	// Decode all frames
	for i := range frames {
		decoded, err := dec.Decode(encoded[i])
		if err != nil {
			t.Fatalf("Decode() frame %d failed: %v", i, err)
		}

		if len(decoded) != len(frames[i]) {
			t.Errorf("Frame %d: decoded length %d, want %d", i, len(decoded), len(frames[i]))
		}
	}
}
