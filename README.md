# go-raknet

**go-raknet** is a Go implementation of the **RakNet protocol**, designed primarily for **San Andreas Multiplayer (SA-MP)** servers.  
It provides reliable UDP communication, packet serialization, and connection management, allowing developers to build custom SA-MP-compatible servers in Go.

> ⚠️ **Note:** This library is **not perfect** and is published as-is because the original maintainer no longer plans to actively develop it. Use at your own discretion.

---

## 🌐 Overview

RakNet is a networking engine widely used in games, offering reliable UDP, packet ordering, and NAT traversal.  
**go-raknet** implements these core features in Go, focusing on SA-MP networking requirements.

This library is suitable for:

- Custom SA-MP server implementations  
- Educational purposes for learning multiplayer networking  
- Experimenting with RakNet protocol in Go  

> **Important:** go-raknet does **not** include SA-MP server binaries or proprietary code. It only provides networking protocol functionality.

---

## ⚡ Features

- Reliable and ordered UDP packet delivery  
- Connection handshake compatible with RakNet protocol  
- Serialization and deserialization of common packet types  
- Event-driven API: `OnConnect`, `OnDisconnect`, `OnPacket`  
- Supports multiple clients and concurrent connections  
- Minimal boilerplate for easy integration into Go projects  

---

## 🟢 Installation
```bash
go get github.com/ventosilenzioso/go-raknet
```

## 🚀 Getting Started

### Creating a Server
```go
package main

import (
    "log"
    "github.com/ventosilenzioso/go-raknet"
)

func main() {
    server, err := raknet.NewServer(":7777")
    if err != nil {
        log.Fatal(err)
    }

    server.OnConnect(func(client *raknet.Client) {
        log.Println("Client connected:", client.Address())
    })

    server.OnDisconnect(func(client *raknet.Client) {
        log.Println("Client disconnected:", client.Address())
    })

    server.OnPacket(func(client *raknet.Client, packet raknet.Packet) {
        log.Printf("Received packet %d from %s\n", packet.ID, client.Address())
    })

    server.Start()
}
```

The event-driven design ensures minimal boilerplate while maintaining full control over connections and packet handling.

## 📚 Architecture

Event-driven design is the core principle of go-raknet:

- `Server` manages client connections and routes packets
- `Client` objects represent individual connected clients
- `Packet` abstraction simplifies RakNet packet serialization and deserialization

This approach allows handling multiple clients concurrently with simple callbacks.

## ⚖️ Legal & Disclaimer

- go-raknet is not affiliated with SA-MP, Rockstar Games, or RakNet.
- This library is intended for educational purposes and experimental networking projects.
- Distribution of official SA-MP binaries or proprietary content is prohibited.
- ⚠️ This library is incomplete and may not cover all RakNet features. Use at your own risk.

# go-raknet

**go-raknet** is a Go implementation of the **RakNet protocol**, designed primarily for **San Andreas Multiplayer (SA-MP)** servers.  
It provides reliable UDP communication, packet serialization, and connection management, allowing developers to build custom SA-MP-compatible servers in Go.

> ⚠️ **Note:** This library is **not perfect** and is published as-is because the original maintainer no longer plans to actively develop it. Use at your own discretion.

---

## 🌐 Overview

RakNet is a networking engine widely used in games, offering reliable UDP, packet ordering, and NAT traversal.  
**go-raknet** implements these core features in Go, focusing on SA-MP networking requirements.

This library is suitable for:

- Custom SA-MP server implementations  
- Educational purposes for learning multiplayer networking  
- Experimenting with RakNet protocol in Go  

> **Important:** go-raknet does **not** include SA-MP server binaries or proprietary code. It only provides networking protocol functionality.

---

## ⚡ Features

- Reliable and ordered UDP packet delivery  
- Connection handshake compatible with RakNet protocol  
- Serialization and deserialization of common packet types  
- Event-driven API: `OnConnect`, `OnDisconnect`, `OnPacket`  
- Supports multiple clients and concurrent connections  
- Minimal boilerplate for easy integration into Go projects  

---

## 🟢 Installation
```bash
go get github.com/ventosilenzioso/go-raknet
```

## 🚀 Getting Started

### Creating a Server
```go
package main

import (
    "log"
    "github.com/ventosilenzioso/go-raknet"
)

func main() {
    server, err := raknet.NewServer(":7777")
    if err != nil {
        log.Fatal(err)
    }

    server.OnConnect(func(client *raknet.Client) {
        log.Println("Client connected:", client.Address())
    })

    server.OnDisconnect(func(client *raknet.Client) {
        log.Println("Client disconnected:", client.Address())
    })

    server.OnPacket(func(client *raknet.Client, packet raknet.Packet) {
        log.Printf("Received packet %d from %s\n", packet.ID, client.Address())
    })

    server.Start()
}
```

The event-driven design ensures minimal boilerplate while maintaining full control over connections and packet handling.

## 📚 Architecture

Event-driven design is the core principle of go-raknet:

- `Server` manages client connections and routes packets
- `Client` objects represent individual connected clients
- `Packet` abstraction simplifies RakNet packet serialization and deserialization

This approach allows handling multiple clients concurrently with simple callbacks.

## ⚖️ Legal & Disclaimer

- go-raknet is not affiliated with SA-MP, Rockstar Games, or RakNet.
- This library is intended for educational purposes and experimental networking projects.
- Distribution of official SA-MP binaries or proprietary content is prohibited.
- ⚠️ This library is incomplete and may not cover all RakNet features. Use at your own risk.

## 📄 Credits

- RakNet protocol: [Facebook RakNet](https://github.com/facebookarchive/RakNet)
- SA-MP protocol research: [SA-MP Wiki](https://wiki.sa-mp.com)
- Community examples and documentation

## 🔗 Resources

- [SA-MP Wiki](https://wiki.sa-mp.com)
- [RakNet Protocol Reference](https://github.com/facebookarchive/RakNet/blob/master/Documentation)
- [Examples folder](./examples) in this repository for usage snippets
