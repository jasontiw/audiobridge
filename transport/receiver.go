package transport

import (
	"fmt"
	"net"
)

// UDPReceiver receives audio packets over UDP
type UDPReceiver struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	buffer []byte
}

// NewUDPReceiver creates a new UDP receiver listening on the specified port
func NewUDPReceiver(port int) (*UDPReceiver, error) {
	addr := &net.UDPAddr{
		IP:   net.IPv4zero, // Listen on all interfaces
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	// Set buffer size (64KB is typical for UDP)
	if err := conn.SetReadBuffer(65536); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to set read buffer: %w", err)
	}

	return &UDPReceiver{
		conn:   conn,
		addr:   addr,
		buffer: make([]byte, 65536), // Max UDP packet size
	}, nil
}

// Receive receives an audio packet
// Returns the packet and the sender address
func (r *UDPReceiver) Receive() (*Packet, *net.UDPAddr, error) {
	n, senderAddr, err := r.conn.ReadFromUDP(r.buffer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to receive packet: %w", err)
	}

	// Parse the packet
	packet := &Packet{}
	if err := packet.Unmarshal(r.buffer[:n]); err != nil {
		// Silently discard invalid packets
		return nil, nil, fmt.Errorf("discarded invalid packet: %w", err)
	}

	return packet, senderAddr, nil
}

// Close closes the UDP connection
func (r *UDPReceiver) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}
