package transport

import (
	"fmt"
	"net"
)

// UDPSender sends audio packets over UDP
type UDPSender struct {
	conn     *net.UDPConn
	addr     *net.UDPAddr
	sequence uint32
}

// NewUDPSender creates a new UDP sender
func NewUDPSender(address string) (*UDPSender, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %w", err)
	}

	return &UDPSender{
		conn:     conn,
		addr:     addr,
		sequence: 0,
	}, nil
}

// Send sends an audio packet to the target address
func (s *UDPSender) Send(payload []byte) error {
	// Create packet with current sequence and timestamp
	packet := NewPacket(s.sequence, uint64(s.sequence)*48000, payload)

	// Marshal to binary
	data, err := packet.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal packet: %w", err)
	}

	// Send via UDP
	_, err = s.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send packet: %w", err)
	}

	// Increment sequence for next packet
	s.sequence++

	return nil
}

// Close closes the UDP connection
func (s *UDPSender) Close() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
