package gamemode

import (
	"log"
	"math/rand"
	"time"
)

// Player represents a connected player
type Player struct {
	ID       uint16
	Name     string
	Score    int
	Money    int
	Skin     int
	Health   float32
	Armour   float32
	Position Vector3
	Rotation float32
	Interior int
	World    int
	Team     int
	Wanted   int
	IsAdmin  bool
	LastSeen time.Time
}

// Vector3 represents 3D coordinates
type Vector3 struct {
	X, Y, Z float32
}

// FreeroamGamemode implements a complex freeroam gamemode
type FreeroamGamemode struct {
	players       map[uint16]*Player
	vehicles      map[uint16]*Vehicle
	spawnPoints   []SpawnPoint
	adminCommands map[string]AdminCommand
	playerCommands map[string]PlayerCommand
}

// SpawnPoint defines a spawn location
type SpawnPoint struct {
	Position Vector3
	Rotation float32
	Skin     int
	Team     int
}

// Vehicle represents a spawned vehicle
type Vehicle struct {
	ID       uint16
	ModelID  int
	Position Vector3
	Rotation float32
	Color1   int
	Color2   int
	Owner    uint16
}

// AdminCommand represents an admin command
type AdminCommand struct {
	Name        string
	Description string
	MinLevel    int
	Handler     func(*Player, []string) string
}

// PlayerCommand represents a player command
type PlayerCommand struct {
	Name        string
	Description string
	Handler     func(*Player, []string) string
}

// NewFreeroamGamemode creates a new freeroam gamemode instance
func NewFreeroamGamemode() *FreeroamGamemode {
	gm := &FreeroamGamemode{
		players:        make(map[uint16]*Player),
		vehicles:       make(map[uint16]*Vehicle),
		spawnPoints:    make([]SpawnPoint, 0),
		adminCommands:  make(map[string]AdminCommand),
		playerCommands: make(map[string]PlayerCommand),
	}
	
	gm.initializeSpawnPoints()
	gm.registerCommands()
	
	return gm
}

// initializeSpawnPoints sets up spawn locations
func (gm *FreeroamGamemode) initializeSpawnPoints() {
	// Los Santos spawns
	gm.spawnPoints = []SpawnPoint{
		{Vector3{1958.3783, 1343.1572, 15.3746}, 270.1425, 0, 0},
		{Vector3{2199.6531, 1393.3678, 10.8203}, 0.0000, 1, 0},
		{Vector3{2483.5977, 1222.8304, 10.8203}, 181.8294, 2, 0},
		{Vector3{2495.0964, 1687.7073, 10.8203}, 0.0000, 3, 0},
		{Vector3{2306.3252, 2442.2158, 10.8203}, 94.3914, 4, 0},
		{Vector3{2197.4092, 2487.7598, 10.8203}, 180.4898, 5, 0},
		{Vector3{1768.2111, 2847.4736, 10.8203}, 270.0000, 6, 0},
		{Vector3{1457.4762, 2773.4868, 10.8203}, 270.0000, 7, 0},
	}
	
	log.Printf("✅ Loaded %d spawn points", len(gm.spawnPoints))
}

// registerCommands registers all available commands
func (gm *FreeroamGamemode) registerCommands() {
	// Player commands
	gm.playerCommands["help"] = PlayerCommand{
		Name:        "help",
		Description: "Show available commands",
		Handler:     gm.cmdHelp,
	}
	
	gm.playerCommands["stats"] = PlayerCommand{
		Name:        "stats",
		Description: "Show your statistics",
		Handler:     gm.cmdStats,
	}
	
	gm.playerCommands["kill"] = PlayerCommand{
		Name:        "kill",
		Description: "Commit suicide",
		Handler:     gm.cmdKill,
	}
	
	gm.playerCommands["v"] = PlayerCommand{
		Name:        "v",
		Description: "Spawn a vehicle",
		Handler:     gm.cmdVehicle,
	}
	
	// Admin commands
	gm.adminCommands["kick"] = AdminCommand{
		Name:        "kick",
		Description: "Kick a player",
		MinLevel:    1,
		Handler:     gm.cmdKick,
	}
	
	gm.adminCommands["ban"] = AdminCommand{
		Name:        "ban",
		Description: "Ban a player",
		MinLevel:    2,
		Handler:     gm.cmdBan,
	}
	
	gm.adminCommands["tp"] = AdminCommand{
		Name:        "tp",
		Description: "Teleport to a player",
		MinLevel:    1,
		Handler:     gm.cmdTeleport,
	}
	
	gm.adminCommands["heal"] = AdminCommand{
		Name:        "heal",
		Description: "Heal a player",
		MinLevel:    1,
		Handler:     gm.cmdHeal,
	}
	
	log.Printf("✅ Registered %d player commands and %d admin commands", 
		len(gm.playerCommands), len(gm.adminCommands))
}

// OnPlayerConnect is called when a player connects
func (gm *FreeroamGamemode) OnPlayerConnect(playerID uint16, name string) {
	player := &Player{
		ID:       playerID,
		Name:     name,
		Score:    0,
		Money:    5000,
		Skin:     0,
		Health:   100.0,
		Armour:   0.0,
		Interior: 0,
		World:    0,
		Team:     0,
		Wanted:   0,
		IsAdmin:  false,
		LastSeen: time.Now(),
	}
	
	gm.players[playerID] = player
	
	log.Printf("🎮 [Gamemode] Player %s (ID: %d) connected", name, playerID)
	gm.SendMessageToAll(0xFFFF00AA, player.Name+" has joined the server")
}

// OnPlayerDisconnect is called when a player disconnects
func (gm *FreeroamGamemode) OnPlayerDisconnect(playerID uint16, reason string) {
	player, exists := gm.players[playerID]
	if !exists {
		return
	}
	
	log.Printf("🎮 [Gamemode] Player %s (ID: %d) disconnected: %s", player.Name, playerID, reason)
	gm.SendMessageToAll(0xFF0000AA, player.Name+" has left the server ("+reason+")")
	
	delete(gm.players, playerID)
}

// OnPlayerSpawn is called when a player spawns
func (gm *FreeroamGamemode) OnPlayerSpawn(playerID uint16) {
	player, exists := gm.players[playerID]
	if !exists {
		return
	}
	
	// Get random spawn point
	spawn := gm.spawnPoints[rand.Intn(len(gm.spawnPoints))]
	
	player.Position = spawn.Position
	player.Rotation = spawn.Rotation
	player.Skin = spawn.Skin
	player.Health = 100.0
	player.Armour = 0.0
	
	log.Printf("🎮 [Gamemode] Player %s spawned at %.2f, %.2f, %.2f", 
		player.Name, spawn.Position.X, spawn.Position.Y, spawn.Position.Z)
	
	gm.SendMessageToPlayer(playerID, 0x00FF00AA, "Welcome to SA-MP Freeroam Server!")
	gm.SendMessageToPlayer(playerID, 0xFFFFFFAA, "Type /help to see available commands")
}

// OnPlayerCommand is called when a player types a command
func (gm *FreeroamGamemode) OnPlayerCommand(playerID uint16, command string, args []string) bool {
	player, exists := gm.players[playerID]
	if !exists {
		return false
	}
	
	// Check player commands
	if cmd, found := gm.playerCommands[command]; found {
		result := cmd.Handler(player, args)
		if result != "" {
			gm.SendMessageToPlayer(playerID, 0xFFFFFFAA, result)
		}
		return true
	}
	
	// Check admin commands
	if cmd, found := gm.adminCommands[command]; found {
		if !player.IsAdmin {
			gm.SendMessageToPlayer(playerID, 0xFF0000AA, "You are not authorized to use this command")
			return true
		}
		
		result := cmd.Handler(player, args)
		if result != "" {
			gm.SendMessageToPlayer(playerID, 0xFFFFFFAA, result)
		}
		return true
	}
	
	return false
}

// Command handlers
func (gm *FreeroamGamemode) cmdHelp(player *Player, args []string) string {
	return "Available commands: /help, /stats, /kill, /v [vehicleid]"
}

func (gm *FreeroamGamemode) cmdStats(player *Player, args []string) string {
	return "Stats - Score: " + string(rune(player.Score)) + " | Money: $" + string(rune(player.Money)) + 
		" | Health: " + string(rune(int(player.Health)))
}

func (gm *FreeroamGamemode) cmdKill(player *Player, args []string) string {
	player.Health = 0.0
	log.Printf("🎮 Player %s committed suicide", player.Name)
	return "You have killed yourself"
}

func (gm *FreeroamGamemode) cmdVehicle(player *Player, args []string) string {
	if len(args) < 1 {
		return "Usage: /v [vehicleid]"
	}
	
	// TODO: Spawn vehicle near player
	return "Vehicle spawned (feature coming soon)"
}

func (gm *FreeroamGamemode) cmdKick(player *Player, args []string) string {
	if len(args) < 1 {
		return "Usage: /kick [playerid]"
	}
	
	// TODO: Kick player
	return "Player kicked (feature coming soon)"
}

func (gm *FreeroamGamemode) cmdBan(player *Player, args []string) string {
	if len(args) < 1 {
		return "Usage: /ban [playerid]"
	}
	
	// TODO: Ban player
	return "Player banned (feature coming soon)"
}

func (gm *FreeroamGamemode) cmdTeleport(player *Player, args []string) string {
	if len(args) < 1 {
		return "Usage: /tp [playerid]"
	}
	
	// TODO: Teleport to player
	return "Teleported (feature coming soon)"
}

func (gm *FreeroamGamemode) cmdHeal(player *Player, args []string) string {
	if len(args) < 1 {
		return "Usage: /heal [playerid]"
	}
	
	// TODO: Heal player
	return "Player healed (feature coming soon)"
}

// SendMessageToPlayer sends a message to a specific player
func (gm *FreeroamGamemode) SendMessageToPlayer(playerID uint16, color uint32, message string) {
	// TODO: Implement actual packet sending
	log.Printf("📨 [To %d] %s", playerID, message)
}

// SendMessageToAll sends a message to all players
func (gm *FreeroamGamemode) SendMessageToAll(color uint32, message string) {
	// TODO: Implement actual packet sending
	log.Printf("📢 [Broadcast] %s", message)
}

// GetPlayer returns a player by ID
func (gm *FreeroamGamemode) GetPlayer(playerID uint16) (*Player, bool) {
	player, exists := gm.players[playerID]
	return player, exists
}

// GetPlayerCount returns the number of connected players
func (gm *FreeroamGamemode) GetPlayerCount() int {
	return len(gm.players)
}
