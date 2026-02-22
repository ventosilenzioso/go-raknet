package server

// SA-MP Packet IDs
const (
	// RakNet packets
	ID_OPEN_CONNECTION_REQUEST  = 0x05
	ID_OPEN_CONNECTION_REPLY    = 0x06
	ID_CONNECTION_REQUEST       = 0x07
	ID_CONNECTION_ACCEPTED      = 0x08
	ID_UNCONNECTED_PING         = 0x1C
	ID_UNCONNECTED_PONG         = 0x1D
	
	// SA-MP specific packets
	ID_PLAYER_SYNC              = 0xCF
	ID_VEHICLE_SYNC             = 0xC8
	ID_PASSENGER_SYNC           = 0xD2
	ID_SPECTATOR_SYNC           = 0xD4
	ID_AIM_SYNC                 = 0xC9
	ID_TRAILER_SYNC             = 0xCA
	ID_UNOCCUPIED_SYNC          = 0xCD
	ID_BULLET_SYNC              = 0xCE
	
	// Server packets
	ID_PLAYER_JOIN              = 0x89
	ID_PLAYER_QUIT              = 0x8A
	ID_SPAWN_PLAYER             = 0x8B
	ID_DEATH_NOTIFICATION       = 0x8C
	ID_SERVER_MESSAGE           = 0x8D
)

type Packet struct {
	ID   byte
	Data []byte
}

func NewPacket(id byte, data []byte) *Packet {
	return &Packet{
		ID:   id,
		Data: data,
	}
}

func (p *Packet) Serialize() []byte {
	result := make([]byte, 1+len(p.Data))
	result[0] = p.ID
	copy(result[1:], p.Data)
	return result
}
