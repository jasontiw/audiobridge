package transport

import (
	"bytes"
	"testing"
)

func TestPacketMarshalUnmarshal(t *testing.T) {
	// Test round-trip of packet marshal/unmarshal
	original := NewPacket(42, 1234567890, []byte{0x01, 0x02, 0x03, 0x04})

	// Marshal to binary
	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal packet: %v", err)
	}

	// Check header size
	if len(data) != HeaderSize+4 { // 12 bytes header + 4 bytes payload
		t.Errorf("Expected packet size %d, got %d", HeaderSize+4, len(data))
	}

	// Unmarshal
	parsed := &Packet{}
	err = parsed.Unmarshal(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal packet: %v", err)
	}

	// Verify fields
	if parsed.Seq != original.Seq {
		t.Errorf("Sequence mismatch: got %d, want %d", parsed.Seq, original.Seq)
	}

	if !bytes.Equal(parsed.Payload, original.Payload) {
		t.Errorf("Payload mismatch: got %v, want %v", parsed.Payload, original.Payload)
	}

	if parsed.GetTimestamp() != original.GetTimestamp() {
		t.Errorf("Timestamp mismatch: got %d, want %d", parsed.GetTimestamp(), original.GetTimestamp())
	}

	// Verify magic
	if string(parsed.Magic[:]) != Magic {
		t.Errorf("Magic mismatch: got %q, want %q", string(parsed.Magic[:]), Magic)
	}

	// Verify version
	if parsed.Version != Version {
		t.Errorf("Version mismatch: got %d, want %d", parsed.Version, Version)
	}
}

func TestPacketInvalidMagic(t *testing.T) {
	// Create a packet with invalid magic
	invalidData := []byte{
		0xFF, 0xFF, // Invalid magic
		0x01,                   // Version
		0x00, 0x00, 0x00, 0x2A, // Sequence (42)
		0x00, 0x00, 0x00, 0x00, 0x00, // Timestamp
		0x01, 0x02, 0x03, 0x04, // Payload
	}

	p := &Packet{}
	err := p.Unmarshal(invalidData)
	if err == nil {
		t.Error("Expected error for invalid magic, got nil")
	}
}

func TestPacketInvalidVersion(t *testing.T) {
	// Create a packet with unsupported version
	invalidData := []byte{
		'A', 'B', // Valid magic
		0xFF,                   // Invalid version
		0x00, 0x00, 0x00, 0x2A, // Sequence
		0x00, 0x00, 0x00, 0x00, 0x00, // Timestamp
		0x01, 0x02, 0x03, 0x04, // Payload
	}

	p := &Packet{}
	err := p.Unmarshal(invalidData)
	if err == nil {
		t.Error("Expected error for invalid version, got nil")
	}
}

func TestPacketTooSmall(t *testing.T) {
	// Packet smaller than header
	smallData := []byte{0x01, 0x02, 0x03}

	p := &Packet{}
	err := p.Unmarshal(smallData)
	if err == nil {
		t.Error("Expected error for small packet, got nil")
	}
}

func TestPacketTimestamp(t *testing.T) {
	// Test 40-bit timestamp handling
	ts := uint64(0x0102030405) // 5-byte timestamp
	original := NewPacket(0, ts, []byte{})

	got := original.GetTimestamp()
	if got != ts {
		t.Errorf("Timestamp mismatch: got 0x%08x, want 0x%08x", got, ts)
	}
}

func TestPacketEmptyPayload(t *testing.T) {
	// Test packet with empty payload
	original := NewPacket(1, 1000, []byte{})

	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal packet: %v", err)
	}

	parsed := &Packet{}
	err = parsed.Unmarshal(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal packet: %v", err)
	}

	if len(parsed.Payload) != 0 {
		t.Errorf("Expected empty payload, got %d bytes", len(parsed.Payload))
	}
}
