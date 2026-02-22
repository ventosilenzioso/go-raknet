package protocol

import (
	"net"
	"testing"
)

func TestBitStreamWriteRead(t *testing.T) {
	bs := NewEmptyBitStream()
	
	// Write data
	bs.WriteByte(0x42)
	bs.WriteUint16(1234)
	bs.WriteUint32(567890)
	bs.WriteString("Hello World")
	
	// Read data
	readBS := NewBitStream(bs.GetData())
	
	b, _ := readBS.ReadByte()
	if b != 0x42 {
		t.Errorf("Expected 0x42, got 0x%02X", b)
	}
	
	u16, _ := readBS.ReadUint16()
	if u16 != 1234 {
		t.Errorf("Expected 1234, got %d", u16)
	}
	
	u32, _ := readBS.ReadUint32()
	if u32 != 567890 {
		t.Errorf("Expected 567890, got %d", u32)
	}
	
	str, _ := readBS.ReadString()
	if str != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", str)
	}
}

func TestEncapsulatedPacket(t *testing.T) {
	packet := &EncapsulatedPacket{
		Reliability:  RELIABLE_ORDERED,
		MessageIndex: 123,
		OrderIndex:   456,
		OrderChannel: 0,
		Split:        false,
		Payload:      []byte{0x01, 0x02, 0x03},
	}
	
	size := packet.GetSize()
	if size <= 0 {
		t.Error("Packet size should be greater than 0")
	}
}

func TestDataPacketEncodeDecode(t *testing.T) {
	dp := NewDataPacket()
	dp.SequenceNumber = 100
	
	encap := &EncapsulatedPacket{
		Reliability: RELIABLE,
		MessageIndex: 50,
		Payload:     []byte{0xAA, 0xBB, 0xCC},
	}
	dp.Packets = append(dp.Packets, encap)
	
	encoded := dp.Encode()
	
	decoded, err := DecodeDataPacket(encoded)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}
	
	if decoded.SequenceNumber != dp.SequenceNumber {
		t.Errorf("Expected sequence %d, got %d", dp.SequenceNumber, decoded.SequenceNumber)
	}
	
	if len(decoded.Packets) != 1 {
		t.Errorf("Expected 1 packet, got %d", len(decoded.Packets))
	}
}

func TestACKEncode(t *testing.T) {
	ack := NewACK()
	ack.Packets = []uint32{1, 2, 3, 4, 5}
	
	encoded := ack.Encode()
	
	if encoded[0] != 0xC0 {
		t.Errorf("Expected ACK flag 0xC0, got 0x%02X", encoded[0])
	}
}

func TestNACKEncode(t *testing.T) {
	nack := NewNACK()
	nack.Packets = []uint32{10, 11, 12}
	
	encoded := nack.Encode()
	
	if encoded[0] != 0xA0 {
		t.Errorf("Expected NACK flag 0xA0, got 0x%02X", encoded[0])
	}
}

func TestSessionCreation(t *testing.T) {
	addr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 7777,
	}
	
	session := NewSession(addr, 1492)
	
	if session.State != STATE_UNCONNECTED {
		t.Errorf("Expected STATE_UNCONNECTED, got %d", session.State)
	}
	
	if session.MTU != 1492 {
		t.Errorf("Expected MTU 1492, got %d", session.MTU)
	}
}

func TestAddressWriteRead(t *testing.T) {
	bs := NewEmptyBitStream()
	
	addr := &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 100),
		Port: 7777,
	}
	
	bs.WriteAddress(addr)
	
	readBS := NewBitStream(bs.GetData())
	readAddr, err := readBS.ReadAddress()
	
	if err != nil {
		t.Fatalf("Failed to read address: %v", err)
	}
	
	if !readAddr.IP.Equal(addr.IP) {
		t.Errorf("Expected IP %s, got %s", addr.IP, readAddr.IP)
	}
	
	if readAddr.Port != addr.Port {
		t.Errorf("Expected port %d, got %d", addr.Port, readAddr.Port)
	}
}
