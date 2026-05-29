# Aerohockey WS Protocol v1

Wire format: JSON over WSS, UTF-8, one message per WS frame, no fragmentation expected.
Every message has a `type` discriminator. All numeric coordinates are in the engine's logical 800×400 coordinate system. `tServer` / `tClient` are Unix ms.

## Client → Server

| type   | payload                                                            | rate    |
|--------|--------------------------------------------------------------------|---------|
| AUTH   | `{ initData: string }`                                             | 1× on connect |
| INPUT  | `{ seq: int, tClient: int64, malletTarget: { x: float, y: float }}`| 30 Hz   |
| PING   | `{ tClient: int64 }`                                               | optional |

## Server → Client

| type        | payload                                                                                              |
|-------------|------------------------------------------------------------------------------------------------------|
| AUTH_OK     | `{ playerSlot: "A"|"B", opponent: { username: string, telegramId: number } }`                        |
| AUTH_FAIL   | `{ reason: string }`                                                                                 |
| MATCH_STATE | `{ phase: "PENDING"|"COUNTDOWN"|"LIVE"|"SETTLED", countdownMs?: int, durationLeftMs?: int }`         |
| SNAPSHOT    | `{ tServer: int64, ackSeq: int, malletA: V2, malletB: V2, puck: { x,y,vx,vy }, score: { a,b } }`     |
| GOAL        | `{ scorer: "A"|"B", score: { a: int, b: int } }`                                                     |
| MATCH_END   | `{ winnerUserId: string|null, winnerSlot: "A"|"B"|null, reason: "score"|"timeout"|"forfeit"|"no_join", finalScore: {a,b} }` |
| PONG        | `{ tClient: int64, tServer: int64 }`                                                                 |

`V2` = `{ x: float, y: float }`.

## Conventions
- Slot **A** defends the left goal, slot **B** the right goal.
- Mallet target positions submitted in `INPUT` are clamped server-side to the player's half.
- `ackSeq` echoes the last applied `INPUT.seq` for client-side reconciliation.

## Changelog
- 2026-05-30 — `MATCH_END` now always broadcast by the server on match end (score cap, timeout, forfeit). Added `winnerSlot` (`"A"` | `"B"` | `null`) alongside `winnerUserId` — clients should prefer `winnerSlot` for win/loss display since they know their own slot but not their server-side UUID.
- 2026-05-12 — v1 initial.
