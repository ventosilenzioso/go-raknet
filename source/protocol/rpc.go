package protocol

import (
	"encoding/binary"
	"math"
)

// RPC IDs for SA-MP 0.3.7
const (
	RPC_InitGame                 = 0x2B // CRITICAL: Must be sent before SetSpawnInfo
	RPC_SetSpawnInfo             = 0x2C
	RPC_SpawnPlayer              = 0x34
	RPC_TogglePlayerControllable = 0x15
	RPC_SetPlayerPos             = 0x0C
	RPC_SetPlayerFacingAngle     = 0x13
	RPC_SetPlayerHealth          = 0x0E
	RPC_SetPlayerArmour          = 0x42
	RPC_GivePlayerWeapon         = 0x16
	RPC_SetPlayerSkin            = 0x99
	RPC_SetGameModeText          = 0x3E // Set gamemode text
	RPC_SetWeather               = 0x0B // Set weather
	RPC_SetWorldTime             = 0x29 // Set world time
	RPC_SetGravity               = 0x92 // Set gravity
)

// Helper functions for little-endian encoding (SA-MP uses little-endian for RPCs)

func writeUint8(buf *[]byte, v uint8) {
	*buf = append(*buf, v)
}

func writeInt32LE(buf *[]byte, v int32) {
	*buf = append(*buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

func writeUint32LE(buf *[]byte, v uint32) {
	*buf = append(*buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

func writeFloat32LE(buf *[]byte, f float32) {
	bits := math.Float32bits(f)
	writeUint32LE(buf, bits)
}

// BuildInitGameRPC builds InitGame RPC payload (0x2B) for SA-MP 0.3.7-R2
// CRITICAL: This MUST be sent before SetSpawnInfo for SA-MP 0.3.7 client
// Structure based on official SA-MP 0.3.7-R2 protocol
func BuildInitGameRPC(
	zoneNames bool,
	useCJWalk bool,
	allowWeapons bool,
	limitGlobalChatRadius bool,
	globalChatRadius float32,
	stuntBonus bool,
	nameTagDrawDistance float32,
	disableEnterExits bool,
	nameTagLOS bool,
	manualVehicleEngineAndLights bool,
	spawnsAvailable uint32,
	playerID uint16,
	showNameTags bool,
	showPlayerMarkers uint32,
	worldTimeHour uint8,
	weather uint8,
	gravity float32,
	lanMode bool,
	deathDropMoney int32,
	instagib bool,
	onFootRate uint32,
	inCarRate uint32,
	weaponRate uint32,
	multiplier uint32,
	lagCompensation uint32,
	hostname string,
	vehicleFriendlyFire bool,
	usePlayerPedAnims bool,
	worldBoundsMinX float32,
	worldBoundsMinY float32,
	worldBoundsMaxX float32,
	worldBoundsMaxY float32,
	gamemodeText string,
	mapName string,
) []byte {
	buf := make([]byte, 0, 512)
	
	// RPC ID
	writeUint8(&buf, RPC_InitGame)
	
	// Zone names enabled
	if zoneNames {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Use CJ walk
	if useCJWalk {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Allow weapons
	if allowWeapons {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Limit global chat radius
	if limitGlobalChatRadius {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Global chat radius
	writeFloat32LE(&buf, globalChatRadius)
	
	// Stunt bonus
	if stuntBonus {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Name tag draw distance
	writeFloat32LE(&buf, nameTagDrawDistance)
	
	// Disable enter/exit markers
	if disableEnterExits {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Name tag LOS
	if nameTagLOS {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Manual vehicle engine and lights
	if manualVehicleEngineAndLights {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Spawns available
	writeUint32LE(&buf, spawnsAvailable)
	
	// Player ID (2 bytes little endian)
	buf = append(buf, byte(playerID), byte(playerID>>8))
	
	// Show name tags
	if showNameTags {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Show player markers
	writeUint32LE(&buf, showPlayerMarkers)
	
	// World time (hour)
	writeUint8(&buf, worldTimeHour)
	
	// Weather
	writeUint8(&buf, weather)
	
	// Gravity
	writeFloat32LE(&buf, gravity)
	
	// LAN mode
	if lanMode {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// Death drop money
	writeInt32LE(&buf, deathDropMoney)
	
	// Instagib
	if instagib {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// OnFoot rate
	writeUint32LE(&buf, onFootRate)
	
	// InCar rate
	writeUint32LE(&buf, inCarRate)
	
	// Weapon rate
	writeUint32LE(&buf, weaponRate)
	
	// Multiplier
	writeUint32LE(&buf, multiplier)
	
	// Lag compensation
	writeUint32LE(&buf, lagCompensation)
	
	// SA-MP 0.3.7-R2: Hostname (string with uint32 length prefix)
	writeUint32LE(&buf, uint32(len(hostname)))
	buf = append(buf, []byte(hostname)...)
	
	// SA-MP 0.3.7-R2: Vehicle friendly fire
	if vehicleFriendlyFire {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// SA-MP 0.3.7-R2: Use player ped anims
	if usePlayerPedAnims {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	
	// SA-MP 0.3.7-R2: World bounds (4 floats)
	writeFloat32LE(&buf, worldBoundsMinX)
	writeFloat32LE(&buf, worldBoundsMinY)
	writeFloat32LE(&buf, worldBoundsMaxX)
	writeFloat32LE(&buf, worldBoundsMaxY)
	
	// SA-MP 0.3.7-R2: Gamemode text (string with uint32 length prefix)
	writeUint32LE(&buf, uint32(len(gamemodeText)))
	buf = append(buf, []byte(gamemodeText)...)
	
	// SA-MP 0.3.7-R2: Map name (string with uint32 length prefix)
	writeUint32LE(&buf, uint32(len(mapName)))
	buf = append(buf, []byte(mapName)...)
	
	return buf
}

// BuildSetSpawnInfoRPC builds SetSpawnInfo RPC payload
// CRITICAL: team MUST be uint8 (1 byte), NOT int32!
func BuildSetSpawnInfoRPC(team uint8, skin int32, x, y, z, rotation float32, weapon1, ammo1, weapon2, ammo2, weapon3, ammo3 int32) []byte {
	buf := make([]byte, 0, 64)
	
	writeUint8(&buf, RPC_SetSpawnInfo)
	writeUint8(&buf, team) // FIXED: team is uint8, not int32!
	writeInt32LE(&buf, skin)
	writeFloat32LE(&buf, x)
	writeFloat32LE(&buf, y)
	writeFloat32LE(&buf, z)
	writeFloat32LE(&buf, rotation)
	
	// Weapon slot 1
	writeInt32LE(&buf, weapon1)
	writeInt32LE(&buf, ammo1)
	
	// Weapon slot 2
	writeInt32LE(&buf, weapon2)
	writeInt32LE(&buf, ammo2)
	
	// Weapon slot 3
	writeInt32LE(&buf, weapon3)
	writeInt32LE(&buf, ammo3)
	
	return buf
}

// BuildSpawnPlayerRPC builds SpawnPlayer RPC payload
func BuildSpawnPlayerRPC() []byte {
	buf := make([]byte, 0, 1)
	writeUint8(&buf, RPC_SpawnPlayer)
	return buf
}

// BuildTogglePlayerControllableRPC builds TogglePlayerControllable RPC payload
func BuildTogglePlayerControllableRPC(toggle bool) []byte {
	buf := make([]byte, 0, 2)
	writeUint8(&buf, RPC_TogglePlayerControllable)
	if toggle {
		writeUint8(&buf, 1)
	} else {
		writeUint8(&buf, 0)
	}
	return buf
}

// BuildSetPlayerPosRPC builds SetPlayerPos RPC payload
func BuildSetPlayerPosRPC(x, y, z float32) []byte {
	buf := make([]byte, 0, 16)
	writeUint8(&buf, RPC_SetPlayerPos)
	writeFloat32LE(&buf, x)
	writeFloat32LE(&buf, y)
	writeFloat32LE(&buf, z)
	return buf
}

// BuildSetPlayerFacingAngleRPC builds SetPlayerFacingAngle RPC payload
func BuildSetPlayerFacingAngleRPC(angle float32) []byte {
	buf := make([]byte, 0, 8)
	writeUint8(&buf, RPC_SetPlayerFacingAngle)
	writeFloat32LE(&buf, angle)
	return buf
}

// EncodeRPCPacket wraps RPC payload with RakNet RPC ID
func EncodeRPCPacket(rpcPayload []byte) []byte {
	// CRITICAL: SA-MP RPC packets start with 0x7C (ID_RPC), NOT 0x19!
	// 0x19 is ID_NEW_INCOMING_CONNECTION (handshake packet)
	// 0x7C is ID_RPC (Remote Procedure Call)
	packet := make([]byte, 0, len(rpcPayload)+1)
	packet = append(packet, 0x7C) // ID_RPC (correct RakNet RPC identifier)
	packet = append(packet, rpcPayload...)
	return packet
}

// Helper to convert uint32 to bytes (little endian)
func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// Helper to convert float32 to bytes (little endian)
func Float32ToBytes(f float32) []byte {
	bits := math.Float32bits(f)
	return Uint32ToBytes(bits)
}

// BuildSetGameModeTextRPC builds SetGameModeText RPC payload (0x3E)
func BuildSetGameModeTextRPC(gamemode string) []byte {
	buf := make([]byte, 0, len(gamemode)+5)
	writeUint8(&buf, RPC_SetGameModeText)
	
	// String length (4 bytes little endian)
	writeUint32LE(&buf, uint32(len(gamemode)))
	
	// String content (no null terminator)
	buf = append(buf, []byte(gamemode)...)
	
	return buf
}

// BuildSetWeatherRPC builds SetWeather RPC payload (0x0B)
func BuildSetWeatherRPC(weather uint8) []byte {
	buf := make([]byte, 0, 2)
	writeUint8(&buf, RPC_SetWeather)
	writeUint8(&buf, weather)
	return buf
}

// BuildSetWorldTimeRPC builds SetWorldTime RPC payload (0x29)
func BuildSetWorldTimeRPC(hour uint8) []byte {
	buf := make([]byte, 0, 2)
	writeUint8(&buf, RPC_SetWorldTime)
	writeUint8(&buf, hour)
	return buf
}

// BuildSetGravityRPC builds SetGravity RPC payload (0x92)
func BuildSetGravityRPC(gravity float32) []byte {
	buf := make([]byte, 0, 5)
	writeUint8(&buf, RPC_SetGravity)
	writeFloat32LE(&buf, gravity)
	return buf
}
