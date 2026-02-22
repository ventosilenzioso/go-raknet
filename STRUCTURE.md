<h1>Project Structure & File Documentation</h1>

<h2>📁 Professional Folder Structure</h2>

<pre>
raknet-go/
├── pkg/                          <i># Reusable library packages</i>
│   ├── logger/                   <i># Colored logging system</i>
│   │   └── logger.go            <i># Logger with [INFO] [WARN] [ERROR] [SUCCESS]</i>
│   └── raknet/                   <i># Core RakNet protocol</i>
│       └── protocol.go          <i># Packet structures, session, encoding/decoding</i>
│
├── source/                       <i># Implementation source code</i>
│   ├── protocol/                 <i># Protocol layer</i>
│   │   ├── raknet.go            <i># Complete RakNet protocol handler</i>
│   │   ├── rpc.go               <i># SA-MP RPC (Remote Procedure Call)</i>
│   │   └── samp_packets.go      <i># SA-MP specific packets</i>
│   └── server/                   <i># Server implementation</i>
│       ├── server.go            <i># Main server (UDP listener, query handler)</i>
│       ├── raknet_handler.go    <i># RakNet packet handler (handshake, spawn, etc)</i>
│       └── player.go            <i># Player management</i>
│
├── core/                         <i># Application layer (SA-MP server)</i>
│   ├── main.go                  <i># Application entry point (with config struct)</i>
│   ├── events/                   <i># Event system</i>
│   │   └── events.go            <i># Event manager (connect, disconnect, spawn, etc)</i>
│   ├── gamemode/                 <i># Gamemode logic</i>
│   │   └── freeroam.go          <i># Freeroam gamemode (commands, spawn points)</i>
│   ├── systems/                  <i># Game systems</i>
│   │   └── vehicle_system.go   <i># Vehicle spawning & management</i>
│   └── commands/                 <i># Command handlers (empty, ready to fill)</i>
│
├── go.mod                        <i># Go module definition</i>
├── go.sum                        <i># Go dependencies checksum</i>
├── LICENSE                       <i># MIT License</i>
├── README.md                     <i># Main documentation (empty, ready to fill)</i>
├── STRUCTURE.md                  <i># Project structure documentation (this file)</i>
└── .gitignore                    <i># Git ignore rules</i>
</pre>

<hr>

<h2>📦 <code>pkg/</code> - Reusable Library Packages</h2>

<h3><code>pkg/logger/logger.go</code></h3>

<p><strong>Purpose:</strong> Colored logging system for better readability</p>

<ul>
<li><code>[INFO]</code> - White - General information</li>
<li><code>[WARN]</code> - Yellow - Warnings</li>
<li><code>[ERROR]</code> - Red - Errors</li>
<li><code>[SUCCESS]</code> - Green - Successful operations</li>
<li><code>[DEBUG]</code> - Gray - Debug information</li>
</ul>

<p><strong>Main Functions:</strong></p>
<ul>
<li><code>logger.Info()</code> - Log information</li>
<li><code>logger.Warn()</code> - Log warnings</li>
<li><code>logger.Error()</code> - Log errors</li>
<li><code>logger.Success()</code> - Log success</li>
<li><code>logger.Fatal()</code> - Log fatal error and exit</li>
<li><code>logger.Banner()</code> - Display startup banner</li>
</ul>

<h3><code>pkg/raknet/protocol.go</code></h3>

<p><strong>Purpose:</strong> Data structures and helper functions for RakNet protocol</p>

<ul>
<li>Packet ID definitions (0x08, 0x1A, 0xA2, 0x19, 0x84-0x8D, 0xC0, 0xA0)</li>
<li>Reliability types (Unreliable, Reliable, Reliable Ordered, etc)</li>
<li>Session states (Unconnected, Connecting, Connected, In-Game)</li>
<li>Helper functions for encoding/decoding (WriteUint24LE, EncodeDatagram, etc)</li>
</ul>

<p><strong>Important Structures:</strong></p>
<ul>
<li><code>Session</code> - Stores client connection state</li>
<li><code>EncapsulatedPacket</code> - Packet with reliability wrapper</li>
<li><code>Datagram</code> - RakNet datagram container</li>
</ul>

<hr>

<h2>🔧 <code>source/</code> - Implementation Source Code</h2>

<h3><code>source/protocol/raknet.go</code></h3>

<p><strong>Purpose:</strong> Complete RakNet protocol handler implementation</p>

<ul>
<li><strong>Handshake:</strong> 0x08 → 0x1A → 0xA2 → 0x19</li>
<li><strong>Session Management:</strong> Tracking sequence, message index, order index</li>
<li><strong>Reliability System:</strong> ACK/NACK, retransmission, packet ordering</li>
<li><strong>Split Packets:</strong> Handling oversized packets</li>
<li><strong>Thread-safe:</strong> Mutex for concurrent access</li>
</ul>

<p><strong>Main Functions:</strong></p>
<ul>
<li><code>NewSession()</code> - Create new session</li>
<li><code>HandleDataPacket()</code> - Process data packet</li>
<li><code>HandleACK()</code> - Process acknowledgment</li>
<li><code>HandleNACK()</code> - Process negative acknowledgment</li>
<li><code>Update()</code> - Send queued packets and ACK/NACK</li>
</ul>

<h3><code>source/protocol/rpc.go</code></h3>

<p><strong>Purpose:</strong> SA-MP RPC (Remote Procedure Call) builder</p>

<ul>
<li><code>BuildInitGameRPC()</code> - RPC 0x2B (game initialization)</li>
<li><code>BuildSetSpawnInfoRPC()</code> - RPC 0x2C (spawn info)</li>
<li><code>BuildSpawnPlayerRPC()</code> - RPC 0x34 (spawn player)</li>
<li><code>BuildSetGameModeTextRPC()</code> - RPC 0x3E (gamemode text)</li>
<li><code>BuildSetWeatherRPC()</code> - RPC 0x0B (weather)</li>
<li><code>BuildSetWorldTimeRPC()</code> - RPC 0x29 (world time)</li>
</ul>

<p><strong>Format:</strong> All RPCs are wrapped with 0x7C (ID_RPC) + RPC ID + payload</p>

<h3><code>source/protocol/samp_packets.go</code></h3>

<p><strong>Purpose:</strong> SA-MP specific packet handlers</p>

<ul>
<li>E3:00 - Challenge request</li>
<li>E3:01 - Challenge response</li>
<li>E3:07 - Spawn list</li>
<li>E3:21 - Player sync</li>
<li>0x88 - Auth packet</li>
<li>0x8A - Join request</li>
</ul>

<h3><code>source/server/server.go</code></h3>

<p><strong>Purpose:</strong> Main server (UDP listener, query handler)</p>

<ul>
<li>UDP socket listener on port 7777</li>
<li>Query handler (SAMP_QUERY_INFO, SAMP_QUERY_RULES, SAMP_QUERY_PLAYERS, SAMP_QUERY_PING)</li>
<li>Session management (create, get, delete)</li>
<li>Server info (name, gamemode, language, weather, worldtime, etc)</li>
</ul>

<p><strong>Main Functions:</strong></p>
<ul>
<li><code>NewServer()</code> - Create server instance</li>
<li><code>Start()</code> - Start UDP listener</li>
<li><code>Stop()</code> - Stop server</li>
<li><code>handleQuery()</code> - Handle SA-MP query packets</li>
</ul>

<h3><code>source/server/raknet_handler.go</code></h3>

<p><strong>Purpose:</strong> RakNet packet handler (handshake, spawn sequence, etc)</p>

<ul>
<li><strong>Handshake Handler:</strong> 0x08, 0x1A, 0xA2, 0x19</li>
<li><strong>Connection Handler:</strong> E3:00, E3:01, 0x22 (login), 0x8A (join)</li>
<li><strong>Spawn Sequence:</strong> InitGame → SetGameModeText → SetWeather → SetWorldTime → SetSpawnInfo → SpawnPlayer</li>
<li><strong>Packet Routing:</strong> Route packets to appropriate handlers</li>
</ul>

<p><strong>Main Functions:</strong></p>
<ul>
<li><code>HandlePacket()</code> - Main packet router</li>
<li><code>sendRakNetDatagram()</code> - Send reliable ordered packet</li>
<li><code>sendSpawnSequence()</code> - Send spawn RPCs</li>
<li><code>handleDataPacket()</code> - Process data packet</li>
</ul>

<h3><code>source/server/player.go</code></h3>

<p><strong>Purpose:</strong> Player management</p>

<ul>
<li>Player data structure</li>
<li>Player state tracking</li>
<li>Player list management</li>
</ul>

<hr>

<h2>🎮 <code>core/</code> - Application Layer (SA-MP Server)</h2>

<h3><code>core/main.go</code></h3>

<p><strong>Purpose:</strong> Application entry point with config struct</p>

<ul>
<li>Config struct for server configuration</li>
<li>Initialize gamemode (Freeroam)</li>
<li>Create server instance</li>
<li>Setup event handlers</li>
<li>Graceful shutdown handling</li>
</ul>

<p><strong>Config Struct:</strong></p>
<pre><code>type Config struct {
    Host       string
    Port       int
    MaxPlayers int
    ServerName string
    GameMode   string
    Language   string
    Weather    int
    WorldTime  int
    MapName    string
    WebURL     string
}</code></pre>

<p><strong>How to Modify Config:</strong></p>
<ul>
<li>Edit directly in <code>loadConfig()</code> function in <code>core/main.go</code></li>
<li>Or load from environment variables</li>
<li>Or create JSON/YAML config file (optional)</li>
</ul>

<h3><code>core/events/events.go</code></h3>

<p><strong>Purpose:</strong> Event system for gamemode</p>

<ul>
<li><strong>Event Types:</strong> PlayerConnect, PlayerDisconnect, PlayerSpawn, PlayerDeath, PlayerCommand, PlayerText, PlayerUpdate, VehicleSpawn, VehicleDestroy</li>
<li><strong>EventManager:</strong> Register and trigger events</li>
<li><strong>EventHandler:</strong> Function that handles events</li>
</ul>

<p><strong>Usage Example:</strong></p>
<pre><code>eventMgr := events.NewEventManager()

// Register handler
eventMgr.Register(events.EventPlayerConnect, func(event events.Event) {
    log.Printf("Player %d connected", event.PlayerID)
})

// Trigger event
eventMgr.Trigger(events.Event{
    Type:     events.EventPlayerConnect,
    PlayerID: 0,
    Data:     "PlayerName",
})</code></pre>

<h3><code>core/gamemode/freeroam.go</code></h3>

<p><strong>Purpose:</strong> Freeroam gamemode logic</p>

<ul>
<li><strong>Player Management:</strong> Connect, disconnect, spawn</li>
<li><strong>Spawn Points:</strong> 8 spawn locations in Los Santos</li>
<li><strong>Commands:</strong> Player commands (/help, /stats, /kill, /v) and admin commands (/kick, /ban, /tp, /heal)</li>
<li><strong>Messaging:</strong> SendMessageToPlayer, SendMessageToAll</li>
</ul>

<p><strong>Structures:</strong></p>
<ul>
<li><code>Player</code> - Player data (ID, name, score, money, health, position, etc)</li>
<li><code>SpawnPoint</code> - Spawn location (position, rotation, skin, team)</li>
<li><code>Vehicle</code> - Vehicle data</li>
<li><code>AdminCommand</code> - Admin command with level requirement</li>
<li><code>PlayerCommand</code> - Regular player command</li>
</ul>

<p><strong>Event Handlers:</strong></p>
<ul>
<li><code>OnPlayerConnect()</code> - When player connects</li>
<li><code>OnPlayerDisconnect()</code> - When player disconnects</li>
<li><code>OnPlayerSpawn()</code> - When player spawns</li>
<li><code>OnPlayerCommand()</code> - When player types command</li>
</ul>

<h3><code>core/systems/vehicle_system.go</code></h3>

<p><strong>Purpose:</strong> Vehicle spawning & management system</p>

<ul>
<li>Spawn vehicle with model ID, position, rotation, colors</li>
<li>Destroy vehicle</li>
<li>Track vehicle owner</li>
<li>Get vehicle data</li>
</ul>

<p><strong>Main Functions:</strong></p>
<ul>
<li><code>SpawnVehicle()</code> - Spawn new vehicle</li>
<li><code>DestroyVehicle()</code> - Remove vehicle</li>
<li><code>GetVehicle()</code> - Get vehicle data</li>
<li><code>GetVehicleCount()</code> - Number of active vehicles</li>
</ul>

<h3><code>core/commands/</code> (Empty)</h3>

<p><strong>Purpose:</strong> Folder for command handlers (ready to fill)</p>

<ul>
<li>Can be filled with more complex command handlers</li>
<li>Separate command logic from gamemode</li>
<li>Example: <code>admin_commands.go</code>, <code>player_commands.go</code>, <code>vehicle_commands.go</code></li>
</ul>

<hr>

<h2>📚 Documentation</h2>

<h3><code>README.md</code></h3>

<p><strong>Purpose:</strong> Main documentation (empty, ready to fill by you)</p>

<p>This is where you should write:</p>
<ul>
<li>Project description and features</li>
<li>Installation instructions</li>
<li>Usage examples</li>
<li>API documentation</li>
<li>Contributing guidelines</li>
<li>License information</li>
</ul>

<h3><code>STRUCTURE.md</code></h3>

<p><strong>Purpose:</strong> Project structure documentation (this file)</p>

<p>Contains detailed explanation of:</p>
<ul>
<li>Folder structure</li>
<li>File purposes and responsibilities</li>
<li>How each component works</li>
<li>Configuration guide</li>
<li>Development notes</li>
</ul>

<hr>

<h2>🔑 Key Points</h2>

<h3>RakNet Layer (100% Complete)</h3>

<ul>
<li>✅ Handshake protocol (0x08, 0x1A, 0xA2, 0x19)</li>
<li>✅ Session management (consistent pointer, no resets)</li>
<li>✅ RakNet counters (monotonic increment)</li>
<li>✅ Packet reliability (Reliable Ordered working)</li>
<li>✅ ACK/NACK format correct</li>
<li>✅ Client communication (0x28 keepalive, 0x8A join, ACKs received)</li>
</ul>

<h3>SA-MP Layer (Working)</h3>

<ul>
<li>✅ Query handler (info, rules, players, ping)</li>
<li>✅ Auth sequence (E3:00, E3:01, 0x22, 0x8A)</li>
<li>✅ Spawn sequence (InitGame, SetGameModeText, SetWeather, SetWorldTime, SetSpawnInfo, SpawnPlayer)</li>
<li>✅ RPC wrapper (0x7C)</li>
</ul>

<h3>Gamemode Layer (Basic)</h3>

<ul>
<li>✅ Event system</li>
<li>✅ Player management</li>
<li>✅ Spawn points</li>
<li>✅ Commands (basic)</li>
<li>⚠️ Vehicle system (basic, needs integration)</li>
<li>⚠️ Advanced features (TODO)</li>
</ul>

<hr>

<h2>🚀 How to Compile & Run</h2>

<pre><code># Compile
go build -o raknet-server.exe ./core

# Run
./raknet-server.exe

# Or run directly
go run ./core</code></pre>

<hr>

<h2>⚙️ How to Modify Configuration</h2>

<p>Edit file <code>core/main.go</code>, function <code>loadConfig()</code>:</p>

<pre><code>func loadConfig() Config {
    return Config{
        Host:       "0.0.0.0",              // Change bind address
        Port:       7777,                    // Change port
        MaxPlayers: 100,                     // Change max players
        ServerName: "My Server",             // Change server name
        GameMode:   "Freeroam v1.0",         // Change gamemode
        Language:   "English",               // Change language
        Weather:    10,                      // Change weather (0-45)
        WorldTime:  12,                      // Change time (0-23)
        MapName:    "San Andreas",           // Change map name
        WebURL:     "github.com/...",        // Change URL
    }
}</code></pre>

<p><strong>Or use Environment Variables:</strong></p>

<pre><code>func loadConfig() Config {
    return Config{
        Host:       getEnv("SERVER_HOST", "0.0.0.0"),
        Port:       getEnvInt("SERVER_PORT", 7777),
        MaxPlayers: getEnvInt("MAX_PLAYERS", 100),
        // ... etc
    }
}</code></pre>

<hr>

<h2>📝 Notes</h2>

<ul>
<li><strong>RakNet protocol</strong> is generic and can be used for other applications</li>
<li><strong>SA-MP specific code</strong> is in <code>source/protocol/rpc.go</code> and <code>source/protocol/samp_packets.go</code></li>
<li><strong>Gamemode</strong> can be replaced with other gamemodes (roleplay, deathmatch, etc)</li>
<li><strong>Event system</strong> makes integration with gamemode easier</li>
<li><strong>Logging</strong> uses colored logger for better readability</li>
<li><strong>Config</strong> uses Go struct, no external file needed</li>
</ul>

<hr>

<p><em>Last Updated: 2026-02-22</em></p>
