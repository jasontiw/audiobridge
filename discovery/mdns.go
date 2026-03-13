package discovery

import (
	"errors"
)

// mDNS discovery stub - full implementation in Phase 3
var (
	ErrNotImplemented = errors.New("mDNS discovery not implemented yet - coming in Phase 3")
)

// Service represents a discovered AudioBridge service
type Service struct {
	Name string
	Host string
	Port int
}

// StartDiscovery starts looking for AudioBridge services on the network
// Returns a channel of discovered services
func StartDiscovery() (<-chan Service, error) {
	// TODO: Implement in Phase 3 with grandcat/zeroconf
	return nil, ErrNotImplemented
}

// StopDiscovery stops the discovery process
func StopDiscovery() error {
	// TODO: Implement in Phase 3
	return ErrNotImplemented
}

// Announce announces this instance on the network for discovery
func Announce(name string, port int) error {
	// TODO: Implement in Phase 3 with grandcat/zeroconf
	return ErrNotImplemented
}

// StopAnnouncement stops announcing this instance
func StopAnnouncement() error {
	// TODO: Implement in Phase 3
	return ErrNotImplemented
}
