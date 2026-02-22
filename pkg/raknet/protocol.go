package raknet

import (
	"encoding/binary"
	"math"
	"sync"
	"time"
)

// RakNet packet IDs
const (
	// Connection packets
	ID_OPEN_CONNECTION_REQUEST_1  = 0x08
	ID_OPEN_CONNECTION_REPLY_1    = 0x1A
	ID_OPEN_CONNECTION_REQUEST_2  = 0xA2
	ID_OPEN_CONNECTION_REPLY_2    = 0x19
	
	// Data packets
	ID_DATA_PACKET_0 = 0x84
	ID_DATA_PACKET_1 = 0x85
	ID_DATA_PACKET_2 = 0x86
	ID_DATA_PACKET_3 = 0x87
	ID_DATA_PACKET_4 = 0x88
	ID_DATA_PACKET_5 = 0x89
	ID_DATA_PACKET_6 = 0x8A
	ID_DATA_PACKET_7 = 0x8B
	ID_DATA_PACKET_8 = 0x8C
	ID_DATA_PACKET_9 = 0x8D
	
	// ACK/NACK
	ID_ACK  = 0xC0
	ID_NACK = 0xA0
)

// Reliability types
const (
	ReliabilityUnreliable = iota
	ReliabilityUnreliableSequenced
	ReliabilityReliable
	ReliabilityReliableOrdered
	ReliabilityReliableSequenced
)

// Session states
const (
	STATE_UNCONNECTED = iota
	STATE_HANDSHAKE_SENT
	STATE_CONNECTING
	STATE_CONNECTED
	STATE_IN_GAME
)

// Default values
const (
	DEFAULT_MTU_SIZE     = 576
	DEFAULT_TIMEOUT      = 30 * time.Second
	MAX_RETRIES          = 5
	ACK_SEND_INTERVAL    = 50 * time.Millisecond
	KEEPALIVE_INTERVAL   = 5 * time.Second
)

// Session represents a RakNet connection session
type Session struct {
	Addr              string
	State             int
	MTU               int
	
	// RakNet counters
	SequenceNumber    uint32
	MessageIndex      uint32
	ChannelOrderIndex map[uint8]uint32
	
	// ACK/NACK queues
	ACKQueue          map[uint32]struct{}
	NACKQueue         map[uint32]struct{}
	
	// Timing
	LastReceiveTime   time.Time
	LastSendTime      time.Time
	
	// Flags
	GameEntrySent     bool
	
	// Mutex for thread safety
	Mu                sync.RWMutex
}

// NewSession creates a new RakNet session
func NewSession(addr string, mtu int) *Session {
	return &Session{
		Addr:              addr,
		State:             STATE_UNCONNECTED,
		MTU:               mtu,
		SequenceNumber:    0,
		MessageIndex:      0,
		ChannelOrderIndex: make(map[uint8]uint32),
		ACKQueue:          make(map[uint32]struct{}),
		NACKQueue:         make(map[uint32]struct{}),
		LastReceiveTime:   time.Now(),
		LastSendTime:      time.Now(),
		GameEntrySent:     false,
	}
}

// EncapsulatedPacket represents an encapsulated RakNet packet
type EncapsulatedPacket struct {
	Reliability   uint8
	MessageIndex  uint32
	OrderIndex    uint32
	OrderChannel  uint8
	Payload       []byte
}

// Datagram represents a RakNet datagram
type Datagram struct {
	SequenceNumber uint32
	Packets        []EncapsulatedPacket
}

// Helper functions for encoding

// WriteUint24LE writes a 24-bit unsigned integer in little-endian
func WriteUint24LE(v uint32) []byte {
	return []byte{byte(v), byte(v >> 8), byte(v >> 16)}
}

// ReadUint24LE reads a 24-bit unsigned integer in little-endian
func ReadUint24LE(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

// WriteUint32LE writes a 32-bit unsigned integer in little-endian
func WriteUint32LE(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// WriteFloat32LE writes a 32-bit float in little-endian
func WriteFloat32LE(f float32) []byte {
	bits := math.Float32bits(f)
	return WriteUint32LE(bits)
}

// EncodeDatagram encodes a datagram into bytes
func EncodeDatagram(seq uint32, packets []EncapsulatedPacket) []byte {
	buf := make([]byte, 0, 1500)
	
	// Packet ID (0x84 for reliable ordered)
	buf = append(buf, ID_DATA_PACKET_0)
	
	// Sequence number (24-bit LE)
	buf = append(buf, WriteUint24LE(seq)...)
	
	// Encode each encapsulated packet
	for _, pkt := range packets {
		// Flags byte
		flags := pkt.Reliability << 5
		buf = append(buf, flags)
		
		// Payload length in bits (16-bit BE)
		lengthBits := uint16(len(pkt.Payload) * 8)
		buf = append(buf, byte(lengthBits>>8), byte(lengthBits))
		
		// Reliable: message index (24-bit LE)
		if pkt.Reliability >= ReliabilityReliable {
			buf = append(buf, WriteUint24LE(pkt.MessageIndex)...)
		}
		
		// Sequenced: sequence index (24-bit LE)
		if pkt.Reliability == ReliabilityUnreliableSequenced || pkt.Reliability == ReliabilityReliableSequenced {
			buf = append(buf, WriteUint24LE(pkt.OrderIndex)...)
		}
		
		// Ordered: order index (24-bit LE) + order channel (1 byte)
		if pkt.Reliability == ReliabilityReliableOrdered {
			buf = append(buf, WriteUint24LE(pkt.OrderIndex)...)
			buf = append(buf, pkt.OrderChannel)
		}
		
		// Payload
		buf = append(buf, pkt.Payload...)
	}
	
	return buf
}

// EncodeACK encodes an ACK packet
func EncodeACK(sequences []uint32) []byte {
	buf := make([]byte, 0, 100)
	buf = append(buf, ID_ACK)
	
	// Record count (16-bit LE)
	count := uint16(len(sequences))
	buf = append(buf, byte(count), byte(count>>8))
	
	// For simplicity, send each sequence as a single record
	for _, seq := range sequences {
		// Single record (not range)
		buf = append(buf, 0x01) // Record type: single
		buf = append(buf, WriteUint24LE(seq)...)
	}
	
	return buf
}

// EncodeNACK encodes a NACK packet
func EncodeNACK(sequences []uint32) []byte {
	buf := make([]byte, 0, 100)
	buf = append(buf, ID_NACK)
	
	// Record count (16-bit LE)
	count := uint16(len(sequences))
	buf = append(buf, byte(count), byte(count>>8))
	
	// For simplicity, send each sequence as a single record
	for _, seq := range sequences {
		// Single record (not range)
		buf = append(buf, 0x01) // Record type: single
		buf = append(buf, WriteUint24LE(seq)...)
	}
	
	return buf
}
