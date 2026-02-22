package protocol

import (
	"net"
	"testing"
)

func BenchmarkBitStreamWrite(b *testing.B) {
	bs := NewEmptyBitStream()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		bs.Reset()
		bs.WriteByte(0x42)
		bs.WriteUint16(1234)
		bs.WriteUint32(567890)
		bs.WriteString("Hello World")
	}
}

func BenchmarkBitStreamRead(b *testing.B) {
	bs := NewEmptyBitStream()
	bs.WriteByte(0x42)
	bs.WriteUint16(1234)
	bs.WriteUint32(567890)
	bs.WriteString("Hello World")
	data := bs.GetData()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		readBS := NewBitStream(data)
		readBS.ReadByte()
		readBS.ReadUint16()
		readBS.ReadUint32()
		readBS.ReadString()
	}
}

func BenchmarkDataPacketEncode(b *testing.B) {
	dp := NewDataPacket()
	dp.SequenceNumber = 100
	
	for i := 0; i < 10; i++ {
		encap := &EncapsulatedPacket{
			Reliability:  RELIABLE_ORDERED,
			MessageIndex: uint32(i),
			OrderIndex:   uint32(i),
			OrderChannel: 0,
			Payload:      make([]byte, 100),
		}
		dp.Packets = append(dp.Packets, encap)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = dp.Encode()
	}
}

func BenchmarkDataPacketDecode(b *testing.B) {
	dp := NewDataPacket()
	dp.SequenceNumber = 100
	
	for i := 0; i < 10; i++ {
		encap := &EncapsulatedPacket{
			Reliability:  RELIABLE_ORDERED,
			MessageIndex: uint32(i),
			OrderIndex:   uint32(i),
			OrderChannel: 0,
			Payload:      make([]byte, 100),
		}
		dp.Packets = append(dp.Packets, encap)
	}
	
	encoded := dp.Encode()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _ = DecodeDataPacket(encoded)
	}
}

func BenchmarkACKEncode(b *testing.B) {
	ack := NewACK()
	for i := uint32(0); i < 100; i++ {
		ack.Packets = append(ack.Packets, i)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = ack.Encode()
	}
}

func BenchmarkSessionAddToQueue(b *testing.B) {
	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7777}
	session := NewSession(addr, 1492)
	
	packet := &EncapsulatedPacket{
		Reliability: RELIABLE_ORDERED,
		Payload:     make([]byte, 100),
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		session.AddToQueue(packet)
	}
}

func BenchmarkAddressWriteRead(b *testing.B) {
	addr := &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 100),
		Port: 7777,
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		bs := NewEmptyBitStream()
		bs.WriteAddress(addr)
		
		readBS := NewBitStream(bs.GetData())
		_, _ = readBS.ReadAddress()
	}
}

func BenchmarkRakNetPacketSerialize(b *testing.B) {
	packet := NewRakNetPacket(0x42)
	packet.Payload = make([]byte, 500)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = packet.Serialize()
	}
}
