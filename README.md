<div align="center">

<h1>📡 RakNet Protocol Specification</h1>

<p>Reliable UDP-based networking protocol for real-time multiplayer games.</p>

<br/>

![Version](https://img.shields.io/badge/version-1.0.0-3b82f6?style=flat-square)
![Protocol](https://img.shields.io/badge/protocol-UDP-10b981?style=flat-square)
![Status](https://img.shields.io/badge/status-stable-10b981?style=flat-square)
![Purpose](https://img.shields.io/badge/purpose-educational-f59e0b?style=flat-square)

</div>

---

## ⚖️ Legal Notice

> **This implementation is created for educational and interoperability purposes only.**

- **Grand Theft Auto: San Andreas** is a trademark of **Rockstar Games, Inc.** and **Take-Two Interactive Software, Inc.** This project is **not affiliated with, endorsed by, or sponsored by** Rockstar Games or Take-Two Interactive in any way.
- **SA-MP (San Andreas Multiplayer)** is a modification developed independently by the SA-MP Team. All SA-MP protocol behavior described herein is derived from public reverse-engineering and community research for compatibility purposes only.
- **RakNet** was originally developed by **Jenkins Software LLC** and later open-sourced. This implementation is based on the open-source BSD-licensed version.
- All trademarks, game titles, and intellectual property referenced remain the sole property of their respective owners.
- This project does **not** distribute any game assets, binaries, or proprietary content.

---

## 🙏 Credits

<table>
  <thead>
    <tr>
      <th>Contribution</th>
      <th>Author / Source</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>RakNet Original Protocol</td>
      <td><a href="http://www.jenkinssoftware.com/">Jenkins Software LLC</a></td>
    </tr>
    <tr>
      <td>SA-MP Protocol Research</td>
      <td><a href="https://sampwiki.blast.hk/">SA-MP Wiki Community</a></td>
    </tr>
    <tr>
      <td>UDP Networking Best Practices</td>
      <td><a href="https://gafferongames.com/post/udp_vs_tcp/">Glenn Fiedler — gafferongames.com</a></td>
    </tr>
    <tr>
      <td>Protocol Reverse Engineering</td>
      <td>Open-source SA-MP community contributors</td>
    </tr>
  </tbody>
</table>

---

## 📋 Overview

RakNet provides the following capabilities over UDP:

| Feature | Description |
|---|---|
| **Reliable Delivery** | Guaranteed packet delivery with retransmission |
| **Ordered Packets** | Maintains packet order per channel |
| **Fragmentation** | Automatic split and reassembly of large payloads |
| **Connection Management** | Full handshake and session lifecycle |
| **Bandwidth Optimization** | Batching, compression, and selective ACK |
| **Low Latency** | Priority queues and immediate send for critical packets |

---

## 🔄 Connection Flow

<details open>
<summary><strong>1. Initial Handshake</strong></summary>

<br/>

```
Client                          Server
  │                               │
  │──── 0x08 ────────────────────▶│  Open Connection Request 1  (4 bytes)
  │                               │
  │◀─── 0x1A ─────────────────────│  Open Connection Reply 1    (Cookie: port XOR)
  │                               │
  │──── 0xA2 ────────────────────▶│  Open Connection Request 2  (4 bytes)
  │                               │
  │◀─── 0x19 ─────────────────────│  Open Connection Reply 2    (Connection accepted)
  │                               │
```

</details>

<details>
<summary><strong>2. SA-MP Authentication</strong> <em>(Optional)</em></summary>

<br/>

```
Client                          Server
  │                               │
  │──── 0x88 ────────────────────▶│  Auth Request
  │                               │
  │◀─── E3:00 ────────────────────│  Challenge      (25 bytes)
  │                               │
  │──── 0x22 ────────────────────▶│  Login Data     (48 bytes)
  │                               │
  │◀─── E3:01 ────────────────────│  Auth Accept    (3 bytes)
  │                               │
```

</details>

<details>
<summary><strong>3. Game Connection</strong></summary>

<br/>

```
Client                          Server
  │                               │
  │──── 0x8A ────────────────────▶│  Join Request
  │                               │
  │◀─── E5 ───────────────────────│  Player Sync
  │◀─── E3:07 ────────────────────│  Spawn List
  │◀─── E3:21 ────────────────────│  Game Entry Complete
  │                               │
  │◀─── Game RPCs ────────────────│  InitGame, SetSpawnInfo, etc.
  │                               │
```

</details>

---

## 📦 Packet Structure

<details open>
<summary><strong>Data Packet</strong> — <code>0x84</code> to <code>0x8D</code></summary>

<br/>

```
┌──────────┬────────────────────────────────┐
│ Packet   │   Sequence Number (24-bit LE)  │
│    ID    │                                │
├──────────┴────────────────────────────────┤
│                                           │
│          Encapsulated Packets...          │
│                                           │
└───────────────────────────────────────────┘
```

</details>

<details>
<summary><strong>Encapsulated Packet</strong></summary>

<br/>

```
┌──────────┬─────────────────────┐
│  Flags   │  Length (16-bit BE) │
├──────────┴─────────────────────┤
│   Message Index (24-bit LE)    │  ← if Reliable
├────────────────────────────────┤
│   Order Index   (24-bit LE)    │  ← if Ordered
├──────────┬─────────────────────┤
│ Channel  │                     │  ← if Ordered
├──────────┘                     │
│         Payload...             │
└────────────────────────────────┘
```

</details>

<details>
<summary><strong>Flags Byte</strong></summary>

<br/>

| Bits | Field | Values |
|---|---|---|
| `7–5` | Reliability Type | `000` Unreliable · `001` Unreliable Seq · `010` Reliable · `011` Reliable Ordered · `100` Reliable Seq |
| `4` | Has Split | `1` = packet is fragmented |
| `3–0` | Reserved | — |

</details>

<details>
<summary><strong>ACK Packet</strong> — <code>0xC0</code> &nbsp;/&nbsp; <strong>NACK Packet</strong> — <code>0xA0</code></summary>

<br/>

```
┌──────────┬─────────────────────┐
│  0xC0    │  Count (16-bit LE)  │
├──────────┼─────────────────────┤
│  Record  │  Sequence (24-bit)  │
│   Type   │                     │
├──────────┴─────────────────────┤
│        ... more records        │
└────────────────────────────────┘
```

> NACK uses the same structure with packet ID `0xA0`.

</details>

---

## 🔁 Reliability Types

| Type | ID | Delivery | Ordering | Best For |
|---|---|---|---|---|
| Unreliable | `0` | ✗ | ✗ | Position updates, voice |
| Unreliable Sequenced | `1` | ✗ | Latest only | Frequent state updates |
| Reliable | `2` | ✓ | ✗ | Important one-off events |
| **Reliable Ordered** | `3` | **✓** | **✓** | **Game events, RPCs** ← most used |
| Reliable Sequenced | `4` | ✓ | Latest only | State synchronization |

---

## 🔌 Session Management

### States

```
UNCONNECTED  ──▶  HANDSHAKE_SENT  ──▶  CONNECTING  ──▶  CONNECTED  ──▶  IN_GAME
```

| State | Description |
|---|---|
| `UNCONNECTED` | No connection established |
| `HANDSHAKE_SENT` | Waiting for handshake completion |
| `CONNECTING` | Connection in progress |
| `CONNECTED` | Connection established |
| `IN_GAME` | Game session active |

### Session Counters

Each session maintains three monotonically increasing counters that **must never reset** during the session lifetime:

| Counter | Increments On |
|---|---|
| `SequenceNumber` | Every datagram sent |
| `MessageIndex` | Every reliable packet |
| `OrderIndex[channel]` | Every ordered packet per channel |

---

## 📐 MTU

Default MTU: **576 bytes** (maximum single UDP packet size)

| Reliability | Safe Payload Size |
|---|---|
| Reliable Ordered | 501 bytes |
| Reliable | 505 bytes |

> Larger payloads are automatically fragmented and reassembled.

---

## ⏱️ Timing

| Parameter | Value |
|---|---|
| ACK Send Interval | `50ms` |
| Keepalive Interval | `5s` |
| Session Timeout | `30s` |
| Retry Delay | `100ms` (exponential backoff) |

**Retransmission triggers:**
1. No ACK received within timeout
2. NACK received from remote
3. After 5 failed retries → disconnection

---

## ⚠️ Error Handling

<details>
<summary><strong>Packet Loss</strong></summary>

- Detected via sequence number gaps
- NACK sent for missing packets
- Automatic retransmission on NACK or timeout

</details>

<details>
<summary><strong>Out-of-Order Packets</strong></summary>

- Packets buffered until predecessor arrives
- `OrderIndex` used for per-channel reordering
- Independent ordering per channel

</details>

<details>
<summary><strong>Duplicate Packets</strong></summary>

- Detected via sequence number comparison
- Silently discarded — no processing
- ACK still sent to remote

</details>

---

## ⚡ Performance

<details>
<summary><strong>Bandwidth Optimization</strong></summary>

| Technique | Description |
|---|---|
| **Batching** | Multiple small packets combined into one datagram |
| **Compression** | Optional payload compression |
| **Selective ACK** | Only ACK received packets |
| **Delayed ACK** | Batch ACKs to reduce header overhead |

</details>

<details>
<summary><strong>Latency Optimization</strong></summary>

| Technique | Description |
|---|---|
| **Immediate Send** | Critical packets bypass batch queue |
| **Priority Queue** | High-priority packets sent first |
| **Congestion Control** | Adaptive send rate based on network conditions |

</details>

---

## 🔒 Security

### Port Obfuscation

The `0x1A` reply packet XOR-encodes the client port to prevent trivial port scanning:

```go
encoded_hi = (port >> 8)    ^ 0x82
encoded_lo = (port & 0xFF)  ^ 0x93
```

### Session Validation

- Each session maintains unique sequence numbers
- Packets from an unknown or mismatched session are rejected
- Inactive sessions time out after **30 seconds**

---

## 🐛 Debugging

### Enable Logging

```go
logger.SetLevel(logger.LevelDebug)
```

### Common Issues

| Issue | Cause | Fix |
|---|---|---|
| Counter Reset | Counters reset mid-session | Ensure counters are never reset |
| Wrong Reliability | Using Unreliable for game events | Use `Reliable Ordered` for most RPCs |
| MTU Exceeded | Payload too large | Split packets or reduce payload size |
| ACK Timeout | High network latency | Tune retry delay and timeout values |

---

## 📚 References

- 📄 [RakNet Documentation](http://www.jenkinssoftware.com/) — Jenkins Software LLC
- 🎮 [SA-MP Protocol Wiki](https://sampwiki.blast.hk/) — SA-MP Community
- 🌐 [UDP vs TCP — Best Practices](https://gafferongames.com/post/udp_vs_tcp/) — Glenn Fiedler

---

## 🗂️ Version History

<details open>
<summary><strong>v1.0.0</strong> — Initial Release</summary>

<br/>

- ✅ Basic RakNet protocol implementation
- ✅ Reliable ordered packet delivery
- ✅ Session lifecycle management
- ✅ ACK / NACK system

</details>

---

<div align="center">

<sub>
This project is not affiliated with, endorsed by, or connected to
<strong>Rockstar Games</strong>, <strong>Take-Two Interactive</strong>, or the <strong>SA-MP Team</strong>.<br/>
All game trademarks and intellectual property belong to their respective owners.<br/>
For educational and interoperability use only.
</sub>

</div>