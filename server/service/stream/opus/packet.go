package opus

// Packet implements the Depacketizer interface for Opus
type Packet struct{}

// Unmarshal parses the RTP payload into an Opus frame
func (p *Packet) Unmarshal(payload []byte) ([]byte, error) {
	return payload, nil
}

// IsPartitionHead checks if this is the head of a partition
func (p *Packet) IsPartitionHead(payload []byte) bool {
	return true
}

// IsPartitionTail checks if this is the tail of a partition
func (p *Packet) IsPartitionTail(marker bool, payload []byte) bool {
	return true
}
