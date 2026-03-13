package jitter

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Buffer holds audio packets to absorb network jitter
type Buffer struct {
	latencyMs   int
	maxFrames   int                      // maximum number of frames to hold
	packets     *list.List               // ordered list of packets
	sequenceMap map[uint32]*list.Element // quick lookup by sequence
	nextSeq     uint32                   // expected next sequence number
	output      chan []byte
	mu          sync.Mutex
	running     bool

	// Stats
	packetsLost    int
	bufferUnderrun int
}

// BufferStats holds jitter buffer statistics
type BufferStats struct {
	CurrentLatency int // Current latency in ms
	PacketLoss     int // Number of packets lost
	BufferUnderrun int // Number of times buffer was empty when needed
}

// NewBuffer creates a new jitter buffer with the specified latency
// latencyMs: target latency in milliseconds (e.g., 50 for 50ms)
func NewBuffer(latencyMs int) (*Buffer, error) {
	if latencyMs < 10 {
		return nil, errors.New("latency must be at least 10ms")
	}
	if latencyMs > 500 {
		return nil, errors.New("latency cannot exceed 500ms")
	}

	// Calculate max frames: latencyMs / 10ms per frame (Opus uses 10ms frames)
	maxFrames := latencyMs / 10

	// Add some headroom for reordering (2x)
	maxFrames = maxFrames * 2

	return &Buffer{
		latencyMs:      latencyMs,
		maxFrames:      maxFrames,
		packets:        list.New(),
		sequenceMap:    make(map[uint32]*list.Element),
		nextSeq:        0,
		output:         make(chan []byte, 10),
		running:        true,
		packetsLost:    0,
		bufferUnderrun: 0,
	}, nil
}

// Push adds a packet to the buffer
// sequence: packet sequence number
// data: opus encoded audio data
func (b *Buffer) Push(sequence uint32, data []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return errors.New("buffer is closed")
	}

	// If this is the first packet, set nextSeq to this packet's sequence
	if b.nextSeq == 0 {
		b.nextSeq = sequence
	}

	// Check if we already have this packet (duplicate)
	if _, exists := b.sequenceMap[sequence]; exists {
		return nil // Ignore duplicate
	}

	// Add to map for quick lookup
	elem := b.packets.PushBack(PacketEntry{
		sequence: sequence,
		data:     data,
	})
	b.sequenceMap[sequence] = elem

	// Remove oldest packets if buffer is full
	for b.packets.Len() > b.maxFrames {
		oldest := b.packets.Front()
		if oldest != nil {
			entry := oldest.Value.(PacketEntry)
			delete(b.sequenceMap, entry.sequence)
			b.packets.Remove(oldest)
		}
	}

	// Try to output packets in order
	b.outputInOrder()

	return nil
}

// outputInOrder outputs packets in sequence order
func (b *Buffer) outputInOrder() {
	for b.packets.Len() > 0 {
		// Check if the oldest packet is the one we expect
		front := b.packets.Front()
		if front == nil {
			break
		}

		entry := front.Value.(PacketEntry)

		// Calculate gap from expected sequence
		gap := int32(entry.sequence - b.nextSeq)

		if gap < 0 {
			// This packet is older than what we're waiting for
			// It might be a late packet, check if we already passed it
			// If nextSeq has advanced too far, discard it
			if b.nextSeq > entry.sequence && (b.nextSeq-entry.sequence) > 10 {
				delete(b.sequenceMap, entry.sequence)
				b.packets.Remove(front)
				continue
			}
			// Otherwise, it's in the right position, output it
		}

		// If gap is too large, we need to wait for missing packets
		// Only output if gap is 0 (exactly what we expect) or we need to fill a small gap
		if gap > int32(b.maxFrames) {
			// Too far ahead, remove old packets to prevent memory bloat
			delete(b.sequenceMap, entry.sequence)
			b.packets.Remove(front)
			b.packetsLost++
			continue
		}

		// If there's a gap, generate silence for missing packets (only if gap is small)
		if gap > 0 && gap <= 5 {
			// We have a small gap, count lost packets and advance nextSeq
			for i := uint32(0); i < uint32(gap); i++ {
				select {
				case b.output <- nil:
					b.packetsLost++
				default:
					// Output channel full
				}
			}
			b.nextSeq = entry.sequence
		}

		// Output this packet
		select {
		case b.output <- entry.data:
			delete(b.sequenceMap, entry.sequence)
			b.packets.Remove(front)
			b.nextSeq++
		default:
			// Output channel full
		}

		break // Process one packet at a time
	}
}

// Pop retrieves the next audio frame from the buffer
// Returns nil for silence (packet loss)
func (b *Buffer) Pop() ([]byte, error) {
	select {
	case data, ok := <-b.output:
		if !ok {
			return nil, errors.New("buffer closed")
		}
		return data, nil
	case <-time.After(10 * time.Millisecond):
		// Timeout - return nil for silence
		b.mu.Lock()
		b.bufferUnderrun++
		b.mu.Unlock()
		return nil, nil
	}
}

// PopBlocking retrieves the next audio frame, blocking until available
func (b *Buffer) PopBlocking() ([]byte, error) {
	data, ok := <-b.output
	if !ok {
		return nil, errors.New("buffer closed")
	}
	return data, nil
}

// Latency returns the configured buffer latency in milliseconds
func (b *Buffer) Latency() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.latencyMs
}

// Stats returns current buffer statistics
func (b *Buffer) Stats() BufferStats {
	b.mu.Lock()
	defer b.mu.Unlock()
	return BufferStats{
		CurrentLatency: b.latencyMs,
		PacketLoss:     b.packetsLost,
		BufferUnderrun: b.bufferUnderrun,
	}
}

// Close closes the buffer and releases resources
func (b *Buffer) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return nil
	}

	b.running = false
	close(b.output)

	// Clear maps and lists
	b.sequenceMap = make(map[uint32]*list.Element)
	b.packets = list.New()

	return nil
}

// PacketEntry holds a single packet in the buffer
type PacketEntry struct {
	sequence uint32
	data     []byte
}
