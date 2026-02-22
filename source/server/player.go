package server

import (
	"net"
	"time"
)

type Player struct {
	ID       int
	Name     string
	Addr     *net.UDPAddr
	Connected bool
	LastPing time.Time
	
	// Game state
	PosX     float32
	PosY     float32
	PosZ     float32
	Angle    float32
	Health   float32
	Armour   float32
	Skin     int
	Interior int
	VirtualWorld int
}

func NewPlayer(id int, addr *net.UDPAddr) *Player {
	return &Player{
		ID:        id,
		Addr:      addr,
		Connected: false,
		LastPing:  time.Now(),
		Health:    100.0,
		Armour:    0.0,
		Skin:      0,
		Interior:  0,
		VirtualWorld: 0,
	}
}

func (p *Player) SetPosition(x, y, z float32) {
	p.PosX = x
	p.PosY = y
	p.PosZ = z
}

func (p *Player) GetPosition() (float32, float32, float32) {
	return p.PosX, p.PosY, p.PosZ
}

func (p *Player) SetHealth(health float32) {
	if health < 0 {
		health = 0
	}
	if health > 100 {
		health = 100
	}
	p.Health = health
}

func (p *Player) IsAlive() bool {
	return p.Health > 0
}
