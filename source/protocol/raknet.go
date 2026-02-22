package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// RakNet protocol constants
const (
	RAKNET_PROTOCOL_VERSION = 9
	MAX_MTU_SIZE           = 1492
	DEFAULT_MTU_SIZE       = 576
	MAX_CHANNELS           = 32
	MAX_SPLIT_PACKET_COUNT = 128
	
	// Safety margin for IP/UDP overhead to prevent IP fragmentation
	// IP header: 20 bytes (or 60 with options)
	// UDP header: 8 bytes
	// Total overhead: 28-68 bytes
	// We use 60 bytes margin to be safe
	MTU_SAFETY_MARGIN = 60
)

// Offline message data ID
var OfflineMessageDataID = []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78}

// RakNet Packet IDs
const (
	ID_CONNECTED_PING                    = 0x00
	ID_UNCONNECTED_PING                  = 0x01
	ID_UNCONNECTED_PING_OPEN_CONNECTIONS = 0x02
	ID_CONNECTED_PONG                    = 0x03
	ID_OPEN_CONNECTION_REQUEST_1         = 0x05
	ID_OPEN_CONNECTION_REPLY_1           = 0x06
	ID_OPEN_CONNECTION_REQUEST_2         = 0x07
	ID_OPEN_CONNECTION_REPLY_2           = 0x08
	ID_CONNECTION_REQUEST                = 0x09
	ID_CONNECTION_REQUEST_ACCEPTED       = 0x10
	ID_NEW_INCOMING_CONNECTION           = 0x13
	ID_DISCONNECTION_NOTIFICATION        = 0x15
	ID_INCOMPATIBLE_PROTOCOL_VERSION     = 0x19
	ID_UNCONNECTED_PONG                  = 0x1C
	ID_ADVERTISE_SYSTEM                  = 0x1D
	ID_RPC                               = 0x7C // RakNet RPC (Remote Procedure Call)
)

// SA-MP Query Packet IDs
const (
	SAMP_QUERY_INFO    = 'i' // Server info
	SAMP_QUERY_RULES   = 'r' // Server rules
	SAMP_QUERY_PLAYERS = 'c' // Client list (detailed)
	SAMP_QUERY_PING    = 'p' // Ping
)

// Reliability types
const (
	UNRELIABLE                = 0
	UNRELIABLE_SEQUENCED      = 1
	RELIABLE                  = 2
	RELIABLE_ORDERED          = 3
	RELIABLE_SEQUENCED        = 4
	UNRELIABLE_WITH_ACK       = 5
	RELIABLE_WITH_ACK         = 6
	RELIABLE_ORDERED_WITH_ACK = 7
)

// Packet priority
const (
	PRIORITY_IMMEDIATE = 0
	PRIORITY_HIGH      = 1
	PRIORITY_MEDIUM    = 2
	PRIORITY_LOW       = 3
)

type BitStream struct {
	data   []byte
	offset int
}

func NewBitStream(data []byte) *BitStream {
	return &BitStream{
		data:   data,
		offset: 0,
	}
}

func NewEmptyBitStream() *BitStream {
	return &BitStream{
		data:   make([]byte, 0),
		offset: 0,
	}
}

func (bs *BitStream) ReadByte() (byte, error) {
	if bs.offset >= len(bs.data) {
		return 0, fmt.Errorf("buffer overflow")
	}
	b := bs.data[bs.offset]
	bs.offset++
	return b, nil
}

func (bs *BitStream) ReadBytes(n int) ([]byte, error) {
	if bs.offset+n > len(bs.data) {
		return nil, fmt.Errorf("buffer overflow")
	}
	result := bs.data[bs.offset : bs.offset+n]
	bs.offset += n
	return result, nil
}

func (bs *BitStream) ReadUint16() (uint16, error) {
	data, err := bs.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

func (bs *BitStream) ReadUint32() (uint32, error) {
	data, err := bs.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

func (bs *BitStream) ReadUint64() (uint64, error) {
	data, err := bs.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(data), nil
}

func (bs *BitStream) ReadString() (string, error) {
	length, err := bs.ReadUint16()
	if err != nil {
		return "", err
	}
	data, err := bs.ReadBytes(int(length))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (bs *BitStream) ReadAddress() (*net.UDPAddr, error) {
	version, err := bs.ReadByte()
	if err != nil {
		return nil, err
	}
	
	var ip net.IP
	if version == 4 {
		ipBytes, err := bs.ReadBytes(4)
		if err != nil {
			return nil, err
		}
		// Invert bytes for IPv4
		for i := range ipBytes {
			ipBytes[i] = ^ipBytes[i]
		}
		ip = net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
	} else {
		return nil, fmt.Errorf("unsupported IP version: %d", version)
	}
	
	port, err := bs.ReadUint16()
	if err != nil {
		return nil, err
	}
	
	return &net.UDPAddr{IP: ip, Port: int(port)}, nil
}

func (bs *BitStream) WriteByte(b byte) {
	bs.data = append(bs.data, b)
}

func (bs *BitStream) WriteBytes(data []byte) {
	bs.data = append(bs.data, data...)
}

func (bs *BitStream) WriteUint16(v uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	bs.data = append(bs.data, buf...)
}

func (bs *BitStream) WriteUint32(v uint32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	bs.data = append(bs.data, buf...)
}

func (bs *BitStream) WriteUint64(v uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	bs.data = append(bs.data, buf...)
}

func (bs *BitStream) WriteString(s string) {
	bs.WriteUint16(uint16(len(s)))
	bs.data = append(bs.data, []byte(s)...)
}

func (bs *BitStream) WriteAddress(addr *net.UDPAddr) {
	if addr.IP.To4() != nil {
		bs.WriteByte(4)
		ip := addr.IP.To4()
		// Invert bytes for IPv4
		for i := 0; i < 4; i++ {
			bs.WriteByte(^ip[i])
		}
		// Port in LITTLE-ENDIAN for SA-MP
		port := uint16(addr.Port)
		bs.WriteByte(byte(port & 0xFF))
		bs.WriteByte(byte((port >> 8) & 0xFF))
	}
}

func (bs *BitStream) GetData() []byte {
	return bs.data
}

func (bs *BitStream) Reset() {
	bs.data = make([]byte, 0)
	bs.offset = 0
}

func (bs *BitStream) Remaining() int {
	return len(bs.data) - bs.offset
}

type RakNetPacket struct {
	PacketID     byte
	Reliability  byte
	MessageIndex uint32
	OrderIndex   uint32
	OrderChannel byte
	Payload      []byte
}

func NewRakNetPacket(id byte) *RakNetPacket {
	return &RakNetPacket{
		PacketID: id,
		Payload:  make([]byte, 0),
	}
}

func (p *RakNetPacket) Serialize() []byte {
	bs := NewEmptyBitStream()
	bs.WriteByte(p.PacketID)
	bs.WriteBytes(p.Payload)
	return bs.GetData()
}

type EncapsulatedPacket struct {
	Reliability  byte
	MessageIndex uint32
	OrderIndex   uint32
	OrderChannel byte
	Split        bool
	SplitCount   uint32
	SplitID      uint16
	SplitIndex   uint32
	Payload      []byte
}

func (ep *EncapsulatedPacket) GetSize() int {
	size := 3 // Flags + length
	if ep.Reliability == RELIABLE || ep.Reliability == RELIABLE_ORDERED || 
	   ep.Reliability == RELIABLE_SEQUENCED || ep.Reliability == RELIABLE_WITH_ACK || 
	   ep.Reliability == RELIABLE_ORDERED_WITH_ACK {
		size += 3 // Message index
	}
	if ep.Reliability == UNRELIABLE_SEQUENCED || ep.Reliability == RELIABLE_SEQUENCED {
		size += 3 // Sequence index
	}
	if ep.Reliability == RELIABLE_ORDERED || ep.Reliability == RELIABLE_ORDERED_WITH_ACK {
		size += 4 // Order index + channel
	}
	if ep.Split {
		size += 10 // Split packet info
	}
	size += len(ep.Payload)
	return size
}

type DataPacket struct {
	SequenceNumber uint32
	Packets        []*EncapsulatedPacket
}

func NewDataPacket() *DataPacket {
	return &DataPacket{
		Packets: make([]*EncapsulatedPacket, 0),
	}
}

func (dp *DataPacket) Encode() []byte {
	bs := NewEmptyBitStream()
	bs.WriteByte(0x80) // Data packet flag
	bs.WriteUint24(dp.SequenceNumber)
	
	for _, packet := range dp.Packets {
		flags := byte(packet.Reliability << 5)
		if packet.Split {
			flags |= 0x10
		}
		bs.WriteByte(flags)
		
		length := uint16(len(packet.Payload) * 8)
		bs.WriteUint16(length)
		
		if packet.Reliability == RELIABLE || packet.Reliability == RELIABLE_ORDERED || 
		   packet.Reliability == RELIABLE_SEQUENCED || packet.Reliability == RELIABLE_WITH_ACK || 
		   packet.Reliability == RELIABLE_ORDERED_WITH_ACK {
			bs.WriteUint24(packet.MessageIndex)
		}
		
		if packet.Reliability == UNRELIABLE_SEQUENCED || packet.Reliability == RELIABLE_SEQUENCED {
			bs.WriteUint24(packet.OrderIndex)
		}
		
		if packet.Reliability == RELIABLE_ORDERED || packet.Reliability == RELIABLE_ORDERED_WITH_ACK {
			bs.WriteUint24(packet.OrderIndex)
			bs.WriteByte(packet.OrderChannel)
		}
		
		if packet.Split {
			bs.WriteUint32(packet.SplitCount)
			bs.WriteUint16(packet.SplitID)
			bs.WriteUint32(packet.SplitIndex)
		}
		
		bs.WriteBytes(packet.Payload)
	}
	
	return bs.GetData()
}

func DecodeDataPacket(data []byte) (*DataPacket, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("packet too short")
	}

	flags := data[0]
	if (flags & 0x80) == 0 {
		return nil, fmt.Errorf("not a data packet")
	}

	// Sequence number (3 bytes, little-endian)
	seqNum := uint32(data[1]) | uint32(data[2])<<8 | uint32(data[3])<<16

	dp := &DataPacket{
		SequenceNumber: seqNum,
		Packets:        make([]*EncapsulatedPacket, 0),
	}

	offset := 4 // Start after flags + seq

	for offset < len(data) {
		if offset >= len(data) {
			break
		}

		packet := &EncapsulatedPacket{}
		encFlags := data[offset]
		offset++

		// ✅ FIX: Reliability = bits[7:5] (3 bits), Split = bit[4]
		packet.Reliability = (encFlags >> 5) & 0x07
		packet.Split = (encFlags & 0x10) != 0

		// ✅ FIX: Length in bits — Big-Endian 16-bit
		if offset+1 >= len(data) {
			break
		}
		lengthBits := uint16(data[offset])<<8 | uint16(data[offset+1])
		offset += 2
		lengthBytes := int((lengthBits + 7) / 8)

		// RELIABLE: has message index (3 bytes)
		if packet.Reliability == RELIABLE ||
			packet.Reliability == RELIABLE_ORDERED ||
			packet.Reliability == RELIABLE_SEQUENCED ||
			packet.Reliability == RELIABLE_WITH_ACK ||
			packet.Reliability == RELIABLE_ORDERED_WITH_ACK {
			if offset+3 > len(data) {
				break
			}
			packet.MessageIndex = uint32(data[offset]) |
				uint32(data[offset+1])<<8 |
				uint32(data[offset+2])<<16
			offset += 3
		}

		// SEQUENCED: has sequence index (3 bytes)
		if packet.Reliability == UNRELIABLE_SEQUENCED ||
			packet.Reliability == RELIABLE_SEQUENCED {
			if offset+3 > len(data) {
				break
			}
			offset += 3 // skip
		}

		// ORDERED: has order index (3 bytes) + channel (1 byte)
		if packet.Reliability == RELIABLE_ORDERED ||
			packet.Reliability == RELIABLE_ORDERED_WITH_ACK {
			if offset+4 > len(data) {
				break
			}
			packet.OrderIndex = uint32(data[offset]) |
				uint32(data[offset+1])<<8 |
				uint32(data[offset+2])<<16
			offset += 3
			packet.OrderChannel = data[offset]
			offset++
		}

		// SPLIT: has split metadata (4+2+4 = 10 bytes)
		if packet.Split {
			if offset+10 > len(data) {
				break
			}
			packet.SplitCount = uint32(data[offset])<<24 |
				uint32(data[offset+1])<<16 |
				uint32(data[offset+2])<<8 |
				uint32(data[offset+3])
			offset += 4
			packet.SplitID = uint16(data[offset])<<8 | uint16(data[offset+1])
			offset += 2
			packet.SplitIndex = uint32(data[offset])<<24 |
				uint32(data[offset+1])<<16 |
				uint32(data[offset+2])<<8 |
				uint32(data[offset+3])
			offset += 4
		}

		// Payload
		if offset+lengthBytes > len(data) {
			break
		}
		packet.Payload = make([]byte, lengthBytes)
		copy(packet.Payload, data[offset:offset+lengthBytes])
		offset += lengthBytes

		dp.Packets = append(dp.Packets, packet)

		// Log inner packet ID for debugging
		if len(packet.Payload) > 0 {
			log.Printf("[DecodeDataPacket] Inner packet ID=0x%02X len=%d reliability=%d",
				packet.Payload[0], len(packet.Payload), packet.Reliability)
		}
	}

	return dp, nil
}

type ACK struct {
	Packets []uint32
}

func NewACK() *ACK {
	return &ACK{
		Packets: make([]uint32, 0),
	}
}

func (ack *ACK) Encode() []byte {
	// CRITICAL: RakNet ACK format (NO extra bytes!)
	// Format: 0xC0 + count(2 bytes LE) + sequences(3 bytes LE each)
	// Example single ACK: C0 01 00 XX XX XX (6 bytes total)
	
	buf := make([]byte, 0, 3+len(ack.Packets)*3)
	
	// Byte 0: ACK ID
	buf = append(buf, 0xC0)
	
	// Bytes 1-2: Record count (little-endian)
	count := uint16(len(ack.Packets))
	buf = append(buf, byte(count))
	buf = append(buf, byte(count>>8))
	
	// Bytes 3+: Sequences (3 bytes little-endian each, NO flag byte!)
	for _, seq := range ack.Packets {
		buf = append(buf, byte(seq))
		buf = append(buf, byte(seq>>8))
		buf = append(buf, byte(seq>>16))
	}
	
	return buf
}

type NACK struct {
	Packets []uint32
}

func NewNACK() *NACK {
	return &NACK{
		Packets: make([]uint32, 0),
	}
}

func (nack *NACK) Encode() []byte {
	// CRITICAL: RakNet NACK format (NO extra bytes!)
	// Format: 0xA0 + count(2 bytes LE) + sequences(3 bytes LE each)
	// Example single NACK: A0 01 00 XX XX XX (6 bytes total)
	
	buf := make([]byte, 0, 3+len(nack.Packets)*3)
	
	// Byte 0: NACK ID
	buf = append(buf, 0xA0)
	
	// Bytes 1-2: Record count (little-endian)
	count := uint16(len(nack.Packets))
	buf = append(buf, byte(count))
	buf = append(buf, byte(count>>8))
	
	// Bytes 3+: Sequences (3 bytes little-endian each, NO flag byte!)
	for _, seq := range nack.Packets {
		buf = append(buf, byte(seq))
		buf = append(buf, byte(seq>>8))
		buf = append(buf, byte(seq>>16))
	}
	
	return buf
}

// Helper functions for uint24
func (bs *BitStream) WriteUint24(v uint32) {
	// RakNet uses 24-bit LITTLE-endian for sequences
	bs.WriteByte(byte(v))
	bs.WriteByte(byte(v >> 8))
	bs.WriteByte(byte(v >> 16))
}

func (bs *BitStream) ReadUint24() (uint32, error) {
	b1, err := bs.ReadByte()
	if err != nil {
		return 0, err
	}
	b2, err := bs.ReadByte()
	if err != nil {
		return 0, err
	}
	b3, err := bs.ReadByte()
	if err != nil {
		return 0, err
	}
	// RakNet uses 24-bit LITTLE-endian for sequences
	return uint32(b1) | uint32(b2)<<8 | uint32(b3)<<16, nil
}

// Helper function to read 24-bit little-endian from byte slice
func ReadUint24LE(b []byte) uint32 {
	if len(b) < 3 {
		return 0
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

// Helper function to write 24-bit little-endian to byte slice
func WriteUint24LE(v uint32) []byte {
	return []byte{
		byte(v & 0xFF),
		byte((v >> 8) & 0xFF),
		byte((v >> 16) & 0xFF),
	}
}

// GetSafePayloadSize - Calculate maximum safe payload size for given MTU
// Returns the maximum payload size that won't cause IP fragmentation
func GetSafePayloadSize(mtu uint16, isOrdered bool) int {
	// Datagram header: 4 bytes (0x84 + 3 byte seq)
	// Encapsulation header for RELIABLE_ORDERED: 11 bytes
	// Encapsulation header for RELIABLE: 7 bytes
	headerSize := 4
	if isOrdered {
		headerSize += 11
	} else {
		headerSize += 7
	}
	
	// Apply safety margin for IP/UDP overhead
	maxSafeSize := int(mtu) - MTU_SAFETY_MARGIN
	
	// Subtract RakNet header
	maxPayload := maxSafeSize - headerSize
	
	if maxPayload < 0 {
		return 0
	}
	
	return maxPayload
}

type Session struct {
	Addr                 *net.UDPAddr
	MTU                  uint16
	GUID                 uint64            // Client GUID for session migration
	
	// Protected by Mu - accessed from multiple goroutines
	State                int
	MessageIndex         uint32
	SequenceNumber       uint32
	OrderIndex           uint32  // DEPRECATED - use ChannelOrderIndex instead
	ChannelOrderIndex    map[uint8]uint32  // Per-channel ordering index (CRITICAL for RakNet)
	SplitID              uint16
	SplitInProgress      bool              // Lock MTU during split packet transmission
	SendQueue            []*EncapsulatedPacket
	RecoveryQueue        map[uint32]*DataPacket
	ACKQueue             map[uint32]struct{}  // Dedup set for ACK sequences
	NACKQueue            []uint32
	SplitPackets         map[uint16]map[uint32]*EncapsulatedPacket
	LastReceiveTime      time.Time
	LastSendTime         time.Time
	LastTenSent          time.Time         // Last time 0x10 was sent (for cooldown)
	Cookie               []byte // SA-MP cookie for session identification
	ReceivedJoinRequest  bool
	HandshakeSent        bool              // Full handshake sequence sent flag
	StreamingDone        bool              // All streaming packets sent flag
	GameEntrySent        bool              // Game entry sequence sent flag
	PostStreamingSent    bool              // Post-streaming sequence sent flag
	JoinResponseSent     bool              // Join response sequence sent flag
	PendingAuth          bool              // Auth packet received, waiting for 0x0B ACK
	AuthSequence         []byte            // Sequence from 0x88
	AuthPayload          []byte            // Payload from 0x88
	PlayerID             uint16            // SA-MP player ID
	Nickname             string            // SA-MP player nickname
	
	// FIX #5: Sent guards to prevent duplicate packets
	SentE3Phase0         bool              // E3:00 challenge sent
	SentE3Phase1         bool              // E3:01 sent
	SentE3Phase7         bool              // E3:07 sent
	SentNWBitStream      bool              // 0xE5 NWBitStream sent
	AuthHandled          bool              // 0x88 auth processed
	
	// Sequence counter for E3 packets
	SendSeq              uint32            // Dynamic sequence for E3 packets (starts at 0, increments with each E3 packet)
	
	Mu                   sync.RWMutex      // Protects all fields above (exported for external access)
	
	// Protected by pendingMu - separate mutex for PendingACK map to avoid deadlock
	PendingACK           map[uint32][]byte // Packets waiting for ACK (for retransmission)
	pendingMu            sync.RWMutex      // Protects PendingACK map
}

const (
	STATE_UNCONNECTED     = 0
	STATE_HANDSHAKE_SENT  = 1
	STATE_CONNECTING      = 2
	STATE_CONNECTED       = 3
	STATE_LOGIN_COMPLETE  = 4  // NEW: Login complete, ready for game entry
	STATE_READY           = 5  // Deprecated - use STATE_LOGIN_COMPLETE
	STATE_IN_GAME         = 6  // Client ready to receive streaming data
)

func NewSession(addr *net.UDPAddr, mtu uint16) *Session {
	s := &Session{
		Addr:              addr,
		MTU:               mtu,
		State:             STATE_UNCONNECTED,
		MessageIndex:      0,
		SequenceNumber:    0,
		OrderIndex:        0,
		ChannelOrderIndex: make(map[uint8]uint32), // Per-channel ordering
		SplitID:           0,
		SendQueue:         make([]*EncapsulatedPacket, 0),
		RecoveryQueue:     make(map[uint32]*DataPacket),
		ACKQueue:          make(map[uint32]struct{}), // Dedup set
		NACKQueue:         make([]uint32, 0),
		SplitPackets:      make(map[uint16]map[uint32]*EncapsulatedPacket),
		PendingACK:        make(map[uint32][]byte),
		LastReceiveTime:   time.Now(),
		LastSendTime:      time.Now(),
	}
	
	// Log safe payload sizes for this MTU
	safeOrdered := GetSafePayloadSize(mtu, true)
	safeReliable := GetSafePayloadSize(mtu, false)
	log.Printf("📊 Session MTU=%d, Safe payload: ORDERED=%d bytes, RELIABLE=%d bytes (margin=%d)", 
		mtu, safeOrdered, safeReliable, MTU_SAFETY_MARGIN)
	
	return s
}

// Thread-safe methods for PendingACK map access
func (s *Session) StorePendingACK(seq uint32, data []byte) {
	s.pendingMu.Lock()
	defer s.pendingMu.Unlock()
	if s.PendingACK == nil {
		s.PendingACK = make(map[uint32][]byte)
	}
	s.PendingACK[seq] = data
}

func (s *Session) GetPendingACK(seq uint32) ([]byte, bool) {
	s.pendingMu.RLock()
	defer s.pendingMu.RUnlock()
	data, exists := s.PendingACK[seq]
	return data, exists
}

func (s *Session) DeletePendingACK(seq uint32) {
	s.pendingMu.Lock()
	defer s.pendingMu.Unlock()
	delete(s.PendingACK, seq)
}

// Thread-safe methods for state flags
func (s *Session) SetHandshakeSent(value bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.HandshakeSent = value
}

func (s *Session) GetHandshakeSent() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.HandshakeSent
}

func (s *Session) SetStreamingDone(value bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.StreamingDone = value
}

func (s *Session) GetStreamingDone() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.StreamingDone
}

func (s *Session) SetGameEntrySent(value bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.GameEntrySent = value
}

func (s *Session) GetGameEntrySent() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.GameEntrySent
}

func (s *Session) SetPostStreamingSent(value bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.PostStreamingSent = value
}

func (s *Session) GetPostStreamingSent() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.PostStreamingSent
}

func (s *Session) SetJoinResponseSent(value bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.JoinResponseSent = value
}

func (s *Session) GetJoinResponseSent() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.JoinResponseSent
}

func (s *Session) UpdateLastReceiveTime() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.LastReceiveTime = time.Now()
}

func (s *Session) GetLastReceiveTime() time.Time {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.LastReceiveTime
}

func (s *Session) AddToQueue(packet *EncapsulatedPacket) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	if packet.Reliability == RELIABLE || packet.Reliability == RELIABLE_ORDERED || 
	   packet.Reliability == RELIABLE_SEQUENCED || packet.Reliability == RELIABLE_WITH_ACK || 
	   packet.Reliability == RELIABLE_ORDERED_WITH_ACK {
		packet.MessageIndex = s.MessageIndex
		s.MessageIndex++
	}
	
	if packet.Reliability == RELIABLE_ORDERED || packet.Reliability == RELIABLE_ORDERED_WITH_ACK {
		packet.OrderIndex = s.OrderIndex
		s.OrderIndex++
	}
	
	s.SendQueue = append(s.SendQueue, packet)
}

func (s *Session) Update(conn *net.UDPConn) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	// FIXED: ACKQueue is now a map (dedup set), convert to slice for sending
	if len(s.ACKQueue) > 0 {
		// Convert map to slice
		ackSeqs := make([]uint32, 0, len(s.ACKQueue))
		for seq := range s.ACKQueue {
			ackSeqs = append(ackSeqs, seq)
		}
		
		if len(ackSeqs) > 0 {
			ack := NewACK()
			ack.Packets = ackSeqs
			ackData := ack.Encode()
			
			// CRITICAL: Verify ACK format
			expectedLen := 3 + len(ackSeqs)*3
			if len(ackData) != expectedLen {
				log.Printf("❌ ERROR: ACK length mismatch! Expected %d, got %d", expectedLen, len(ackData))
			}
			
			n, err := conn.WriteToUDP(ackData, s.Addr)
			if err != nil {
				log.Printf("❌ Failed to send ACK: %v", err)
			} else {
				log.Printf("✅ Sent ACK to %s: %d bytes, %d sequences (deduped)", s.Addr.String(), n, len(ackSeqs))
				log.Printf("   ACK sequences: %v", ackSeqs)
				log.Printf("   ACK hex: %02X", ackData)
				
				// Verify format for single ACK
				if len(ackSeqs) == 1 {
					if len(ackData) == 6 {
						log.Printf("   ✅ ACK format CORRECT: 6 bytes for single record")
						log.Printf("   Format: [0]=0xC0 [1-2]=count(LE)=0x%02X%02X [3-5]=seq(LE)=0x%02X%02X%02X", 
							ackData[1], ackData[2], ackData[3], ackData[4], ackData[5])
					} else {
						log.Printf("   ❌ ACK format WRONG: %d bytes (expected 6 for single record)", len(ackData))
					}
				}
			}
		}
		
		// Clear ACK queue (recreate map)
		s.ACKQueue = make(map[uint32]struct{})
	}
	
	// Send NACKs
	if len(s.NACKQueue) > 0 {
		nack := NewNACK()
		nack.Packets = s.NACKQueue
		conn.WriteToUDP(nack.Encode(), s.Addr)
		s.NACKQueue = make([]uint32, 0)
	}
	
	// Send queued packets
	if len(s.SendQueue) > 0 {
		dp := NewDataPacket()
		dp.SequenceNumber = s.SequenceNumber
		s.SequenceNumber++
		
		for len(s.SendQueue) > 0 && len(dp.Packets) < 120 {
			packet := s.SendQueue[0]
			s.SendQueue = s.SendQueue[1:]
			dp.Packets = append(dp.Packets, packet)
		}
		
		data := dp.Encode()
		n, err := conn.WriteToUDP(data, s.Addr)
		if err != nil {
			log.Printf("❌ Failed to send data packet: %v", err)
		} else {
			log.Printf("📤 Sent data packet to %s: %d bytes, seq: %d, encap packets: %d", 
				s.Addr.String(), n, dp.SequenceNumber, len(dp.Packets))
			log.Printf("   Data packet hex (first 64 bytes): %x", data[:min(64, len(data))])
		}
		s.RecoveryQueue[dp.SequenceNumber] = dp
		s.LastSendTime = time.Now()
	}
	
	return nil
}

func (s *Session) HandleDataPacket(dp *DataPacket) []*RakNetPacket {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	// CRITICAL: Don't add empty packets to ACK queue (SA-MP behavior)
	if len(dp.Packets) > 0 {
		s.ACKQueue[dp.SequenceNumber] = struct{}{} // Dedup set
	}
	s.LastReceiveTime = time.Now()
	
	packets := make([]*RakNetPacket, 0)
	
	for _, encap := range dp.Packets {
		// CRITICAL: Process reliable ordered state machine
		if encap.Reliability == RELIABLE_ORDERED || encap.Reliability == RELIABLE_ORDERED_WITH_ACK {
			// Check if this is a duplicate or out-of-order message
			channel := encap.OrderChannel
			
			// Initialize expected ordering index for this channel if needed
			if s.ChannelOrderIndex == nil {
				s.ChannelOrderIndex = make(map[uint8]uint32)
			}
			
			expectedOrderIndex := s.ChannelOrderIndex[channel]
			
			// DUPLICATE DETECTION: If order index < expected, this is a duplicate
			if encap.OrderIndex < expectedOrderIndex {
				log.Printf("🔄 DUPLICATE: Received order=%d, expected=%d (channel=%d) - IGNORING", 
					encap.OrderIndex, expectedOrderIndex, channel)
				continue // Skip duplicate
			}
			
			// OUT-OF-ORDER DETECTION: If order index > expected, buffer it
			if encap.OrderIndex > expectedOrderIndex {
				log.Printf("⏸️ OUT-OF-ORDER: Received order=%d, expected=%d (channel=%d) - BUFFERING", 
					encap.OrderIndex, expectedOrderIndex, channel)
				// TODO: Implement out-of-order buffering if needed
				// For now, we'll process it anyway (SA-MP might not need strict ordering)
			}
			
			// IN-ORDER: Process this message and update expected index
			if encap.OrderIndex == expectedOrderIndex {
				log.Printf("✅ IN-ORDER: Received order=%d (channel=%d) - PROCESSING", 
					encap.OrderIndex, channel)
				s.ChannelOrderIndex[channel] = expectedOrderIndex + 1
			}
		}
		
		// Process split packets
		if encap.Split {
			if _, exists := s.SplitPackets[encap.SplitID]; !exists {
				s.SplitPackets[encap.SplitID] = make(map[uint32]*EncapsulatedPacket)
			}
			s.SplitPackets[encap.SplitID][encap.SplitIndex] = encap
			
			if uint32(len(s.SplitPackets[encap.SplitID])) == encap.SplitCount {
				var buffer bytes.Buffer
				for i := uint32(0); i < encap.SplitCount; i++ {
					buffer.Write(s.SplitPackets[encap.SplitID][i].Payload)
				}
				delete(s.SplitPackets, encap.SplitID)
				
				if len(buffer.Bytes()) > 0 {
					packet := &RakNetPacket{
						PacketID: buffer.Bytes()[0],
						Payload:  buffer.Bytes()[1:],
					}
					packets = append(packets, packet)
				}
			}
		} else {
			if len(encap.Payload) > 0 {
				packet := &RakNetPacket{
					PacketID: encap.Payload[0],
					Payload:  encap.Payload[1:],
				}
				packets = append(packets, packet)
			}
		}
	}
	
	return packets
}

func (s *Session) HandleACK(data []byte) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	bs := NewBitStream(data)
	bs.ReadByte() // Skip flag
	
	count, _ := bs.ReadUint16()
	for i := uint16(0); i < count; i++ {
		bs.ReadByte() // Skip single/range flag
		start, _ := bs.ReadUint24()
		end, _ := bs.ReadUint24()
		
		for seq := start; seq <= end; seq++ {
			delete(s.RecoveryQueue, seq)
		}
	}
}

func (s *Session) HandleNACK(data []byte) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	bs := NewBitStream(data)
	bs.ReadByte() // Skip flag
	
	count, _ := bs.ReadUint16()
	for i := uint16(0); i < count; i++ {
		bs.ReadByte() // Skip single/range flag
		start, _ := bs.ReadUint24()
		end, _ := bs.ReadUint24()
		
		for seq := start; seq <= end; seq++ {
			if dp, exists := s.RecoveryQueue[seq]; exists {
				for _, packet := range dp.Packets {
					s.SendQueue = append(s.SendQueue, packet)
				}
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CanStream - Check if session is ready to receive streaming data
func (s *Session) CanStream() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.State == STATE_IN_GAME
}

// StopStreaming - Stop all streaming and reset streaming flags
func (s *Session) StopStreaming() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	s.HandshakeSent = false
	s.StreamingDone = false
	s.GameEntrySent = false
	s.PostStreamingSent = false
	s.JoinResponseSent = false
	
	// Clear send queue to stop pending transmissions
	s.SendQueue = nil
}

// NextSeq increments and returns the next E3 packet sequence (3 bytes, little-endian)
func (s *Session) NextSeq() []byte {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.SendSeq++
	return []byte{
		byte(s.SendSeq),
		byte(s.SendSeq >> 8),
		byte(s.SendSeq >> 16),
	}
}
