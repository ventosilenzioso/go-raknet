package server

import (
	"fmt"
	"log"
	"net"
	"samp-server-go/source/protocol"
	"sync"
	"time"
)

type Server struct {
	Host          string
	Port          int
	MaxPlayers    int
	ServerName    string
	GameMode      string
	Language      string
	Weather       int
	WorldTime     int
	MapName       string
	WebURL        string
	Players       map[int]*Player
	conn          *net.UDPConn
	raknet        *RakNetHandler
	mu            sync.RWMutex
	running       bool
	nextPlayerID  int
}

func NewServer(host string, port int, maxPlayers int) *Server {
	return &Server{
		Host:         host,
		Port:         port,
		MaxPlayers:   maxPlayers,
		ServerName:   "SA-MP Server in Go",
		GameMode:     "Freeroam",
		Language:     "English",
		Weather:      10,
		WorldTime:    12,
		MapName:      "San Andreas",
		WebURL:       "www.sa-mp.com",
		Players:      make(map[int]*Player),
		running:      false,
		nextPlayerID: 0,
	}
}

func (s *Server) Start() error {
	addr := &net.UDPAddr{
		IP:   net.ParseIP(s.Host),
		Port: s.Port,
	}
	
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind UDP socket: %w", err)
	}
	
	s.conn = conn
	s.raknet = NewRakNetHandler(conn, s)
	s.running = true
	
	// Set packet handler
	s.raknet.SetPacketHandler(s.handleGamePacket)
	
	log.Printf("Server started on %s:%d", s.Host, s.Port)
	log.Printf("Server Name: %s", s.ServerName)
	log.Printf("Game Mode: %s", s.GameMode)
	log.Printf("Max Players: %d", s.MaxPlayers)
	
	// Start update ticker
	go s.updateLoop()
	
	// Start session cleanup ticker (every 5 seconds)
	go s.sessionCleanupLoop()
	
	return s.listen()
}

func (s *Server) listen() error {
	buffer := make([]byte, 2048)
	
	log.Printf("Listening for packets on %s:%d...", s.Host, s.Port)
	
	for s.running {
		n, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			if s.running {
				log.Printf("Error reading UDP packet: %v", err)
			}
			continue
		}
		
		// Make a copy of the data
		data := make([]byte, n)
		copy(data, buffer[:n])
		
		// Log first byte of every packet for debugging
		if len(data) > 0 && data[0] != 'S' { // Don't log SAMP queries
			log.Printf("Raw packet: 0x%02X (%d bytes) from %s", data[0], n, addr.String())
		}
		
		go s.raknet.HandlePacket(data, addr)
	}
	
	return nil
}

func (s *Server) updateLoop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	
	for s.running {
		<-ticker.C
		s.raknet.Update()
	}
}

// sessionCleanupLoop - Clean up stale sessions based on REAL timeout
func (s *Server) sessionCleanupLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for s.running {
		<-ticker.C
		s.raknet.CleanupStaleSessions()
	}
}

func (s *Server) handleGamePacket(session *protocol.Session, packet *protocol.RakNetPacket) {
	switch packet.PacketID {
	case 0x25: // ID_AUTH_KEY - SA-MP client authentication
		s.handleAuthKey(session, packet)
	case ID_PLAYER_JOIN:
		s.handlePlayerJoin(session, packet)
	case ID_PLAYER_SYNC:
		s.handlePlayerSync(session, packet)
	case ID_VEHICLE_SYNC:
		s.handleVehicleSync(session, packet)
	case ID_SPAWN_PLAYER:
		s.handleSpawnPlayer(session, packet)
	default:
		log.Printf("Unhandled game packet: 0x%02X from %s", packet.PacketID, session.Addr.String())
	}
}

func (s *Server) handleAuthKey(session *protocol.Session, packet *protocol.RakNetPacket) {
	log.Printf("Received AUTH_KEY (0x25) from %s", session.Addr.String())
	log.Printf("Auth key payload length: %d bytes", len(packet.Payload))
	
	// SA-MP client sends auth key after connection established
	// Server should acknowledge and allow client to proceed
	session.State = protocol.STATE_READY
	log.Printf("Client %s authenticated and ready", session.Addr.String())
}

func (s *Server) handlePlayerJoin(session *protocol.Session, packet *protocol.RakNetPacket) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if len(s.Players) >= s.MaxPlayers {
		log.Printf("Server full, rejecting player from %s", session.Addr.String())
		return
	}
	
	playerID := s.nextPlayerID
	s.nextPlayerID++
	
	player := NewPlayer(playerID, session.Addr)
	player.Connected = true
	s.Players[playerID] = player
	
	log.Printf("Player %d joined from %s", playerID, session.Addr.String())
	
	// Send welcome message
	s.sendServerMessage(session, fmt.Sprintf("Welcome to %s!", s.ServerName))
}

func (s *Server) handlePlayerSync(session *protocol.Session, packet *protocol.RakNetPacket) {
	// Handle player position sync
	// This would parse position data and update player state
}

func (s *Server) handleVehicleSync(session *protocol.Session, packet *protocol.RakNetPacket) {
	// Handle vehicle sync
}

func (s *Server) handleSpawnPlayer(session *protocol.Session, packet *protocol.RakNetPacket) {
	// Handle player spawn
	log.Printf("Player spawned from %s", session.Addr.String())
}

func (s *Server) sendServerMessage(session *protocol.Session, message string) {
	response := protocol.NewEmptyBitStream()
	response.WriteByte(ID_SERVER_MESSAGE)
	response.WriteString(message)
	
	packet := &protocol.RakNetPacket{
		PacketID: ID_SERVER_MESSAGE,
		Payload:  response.GetData()[1:],
	}
	
	s.raknet.SendPacket(session, packet, protocol.RELIABLE_ORDERED)
}

func (s *Server) GetPlayerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.Players)
}

func (s *Server) BroadcastMessage(message string) {
	sessions := s.raknet.GetSessions()
	for _, session := range sessions {
		if session.State == protocol.STATE_CONNECTED {
			s.sendServerMessage(session, message)
		}
	}
}

func (s *Server) Stop() {
	log.Println("Stopping server...")
	s.running = false
	
	if s.conn != nil {
		s.conn.Close()
	}
	
	log.Println("Server stopped")
}
