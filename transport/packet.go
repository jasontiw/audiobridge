package transport

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Packet constants
const (
	Magic      = "AB"
	Version    = 0x01
	HeaderSize = 12
)

// Packet represents an audio packet for UDP transmission
type Packet struct {
	Magic     [2]byte
	Version   uint8
	Seq       uint32
	Timestamp [5]byte // 40-bit timestamp
	Payload   []byte
}

// NewPacket creates a new packet with the given sequence number and payload
func NewPacket(seq uint32, timestamp uint64, payload []byte) *Packet {
	p := &Packet{
		Version: Version,
		Seq:     seq,
		Payload: payload,
	}

	// Set magic
	copy(p.Magic[:], Magic)

	// Set 40-bit timestamp (lower 40 bits of 64-bit timestamp)
	p.Timestamp[0] = byte(timestamp >> 32)
	p.Timestamp[1] = byte(timestamp >> 24)
	p.Timestamp[2] = byte(timestamp >> 16)
	p.Timestamp[3] = byte(timestamp >> 8)
	p.Timestamp[4] = byte(timestamp)

	return p
}

// Marshal serializes the packet to binary
func (p *Packet) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write header
	if err := binary.Write(buf, binary.BigEndian, p.Magic[:]); err != nil {
		return nil, fmt.Errorf("failed to write magic: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Version); err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Seq); err != nil {
		return nil, fmt.Errorf("failed to write sequence: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Timestamp[:]); err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Write payload
	if _, err := buf.Write(p.Payload); err != nil {
		return nil, fmt.Errorf("failed to write payload: %w", err)
	}

	return buf.Bytes(), nil
}

// Unmarshal deserializes a packet from binary
func (p *Packet) Unmarshal(data []byte) error {
	if len(data) < HeaderSize {
		return fmt.Errorf("packet too small: %d bytes, need at least %d", len(data), HeaderSize)
	}

	buf := bytes.NewReader(data)

	// Read magic
	if err := binary.Read(buf, binary.BigEndian, p.Magic[:]); err != nil {
		return fmt.Errorf("failed to read magic: %w", err)
	}

	// Verify magic
	if string(p.Magic[:]) != Magic {
		return fmt.Errorf("invalid magic: %q", string(p.Magic[:]))
	}

	// Read version
	if err := binary.Read(buf, binary.BigEndian, &p.Version); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}

	if p.Version != Version {
		return fmt.Errorf("unsupported version: %d", p.Version)
	}

	// Read sequence
	if err := binary.Read(buf, binary.BigEndian, &p.Seq); err != nil {
		return fmt.Errorf("failed to read sequence: %w", err)
	}

	// Read timestamp
	if err := binary.Read(buf, binary.BigEndian, p.Timestamp[:]); err != nil {
		return fmt.Errorf("failed to read timestamp: %w", err)
	}

	// Read payload
	p.Payload = data[HeaderSize:]

	return nil
}

// GetTimestamp returns the 40-bit timestamp as a uint64
func (p *Packet) GetTimestamp() uint64 {
	return uint64(p.Timestamp[0])<<32 |
		uint64(p.Timestamp[1])<<24 |
		uint64(p.Timestamp[2])<<16 |
		uint64(p.Timestamp[3])<<8 |
		uint64(p.Timestamp[4])
}
