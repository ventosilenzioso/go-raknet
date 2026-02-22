package protocol

import (
	"testing"
)

func TestACKEncodeSingleRecord(t *testing.T) {
	ack := NewACK()
	ack.Packets = []uint32{0x123456}
	
	data := ack.Encode()
	
	// Verify length: 3 bytes header + 3 bytes sequence = 6 bytes total
	if len(data) != 6 {
		t.Errorf("ACK length = %d, want 6", len(data))
	}
	
	// Verify format
	if data[0] != 0xC0 {
		t.Errorf("ACK ID = 0x%02X, want 0xC0", data[0])
	}
	
	// Verify count (little-endian)
	count := uint16(data[1]) | uint16(data[2])<<8
	if count != 1 {
		t.Errorf("ACK count = %d, want 1", count)
	}
	
	// Verify sequence (little-endian 24-bit)
	seq := uint32(data[3]) | uint32(data[4])<<8 | uint32(data[5])<<16
	if seq != 0x123456 {
		t.Errorf("ACK sequence = 0x%06X, want 0x123456", seq)
	}
	
	// Verify exact hex format
	expected := []byte{0xC0, 0x01, 0x00, 0x56, 0x34, 0x12}
	for i := 0; i < len(expected); i++ {
		if data[i] != expected[i] {
			t.Errorf("ACK[%d] = 0x%02X, want 0x%02X", i, data[i], expected[i])
		}
	}
}

func TestACKEncodeMultipleRecords(t *testing.T) {
	ack := NewACK()
	ack.Packets = []uint32{0x000001, 0x000002, 0x000003}
	
	data := ack.Encode()
	
	// Verify length: 3 bytes header + 3*3 bytes sequences = 12 bytes total
	expectedLen := 3 + 3*3
	if len(data) != expectedLen {
		t.Errorf("ACK length = %d, want %d", len(data), expectedLen)
	}
	
	// Verify ACK ID
	if data[0] != 0xC0 {
		t.Errorf("ACK ID = 0x%02X, want 0xC0", data[0])
	}
	
	// Verify count
	count := uint16(data[1]) | uint16(data[2])<<8
	if count != 3 {
		t.Errorf("ACK count = %d, want 3", count)
	}
	
	// Verify sequences
	for i := 0; i < 3; i++ {
		offset := 3 + i*3
		seq := uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16
		expected := uint32(i + 1)
		if seq != expected {
			t.Errorf("ACK sequence[%d] = %d, want %d", i, seq, expected)
		}
	}
}

func TestNACKEncodeSingleRecord(t *testing.T) {
	nack := NewNACK()
	nack.Packets = []uint32{0xABCDEF}
	
	data := nack.Encode()
	
	// Verify length: 3 bytes header + 3 bytes sequence = 6 bytes total
	if len(data) != 6 {
		t.Errorf("NACK length = %d, want 6", len(data))
	}
	
	// Verify format
	if data[0] != 0xA0 {
		t.Errorf("NACK ID = 0x%02X, want 0xA0", data[0])
	}
	
	// Verify count (little-endian)
	count := uint16(data[1]) | uint16(data[2])<<8
	if count != 1 {
		t.Errorf("NACK count = %d, want 1", count)
	}
	
	// Verify sequence (little-endian 24-bit)
	seq := uint32(data[3]) | uint32(data[4])<<8 | uint32(data[5])<<16
	if seq != 0xABCDEF {
		t.Errorf("NACK sequence = 0x%06X, want 0xABCDEF", seq)
	}
	
	// Verify exact hex format (0xABCDEF in little-endian = EF CD AB)
	expected := []byte{0xA0, 0x01, 0x00, 0xEF, 0xCD, 0xAB}
	for i := 0; i < len(expected); i++ {
		if data[i] != expected[i] {
			t.Errorf("NACK[%d] = 0x%02X, want 0x%02X", i, data[i], expected[i])
		}
	}
}

func TestACKEncodeEmpty(t *testing.T) {
	ack := NewACK()
	ack.Packets = []uint32{}
	
	data := ack.Encode()
	
	// Verify length: 3 bytes header only
	if len(data) != 3 {
		t.Errorf("Empty ACK length = %d, want 3", len(data))
	}
	
	// Verify format
	if data[0] != 0xC0 {
		t.Errorf("ACK ID = 0x%02X, want 0xC0", data[0])
	}
	
	// Verify count is 0
	count := uint16(data[1]) | uint16(data[2])<<8
	if count != 0 {
		t.Errorf("Empty ACK count = %d, want 0", count)
	}
}
