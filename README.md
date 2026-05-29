# crypto-hockey-game-engine

Go service that runs authoritative aerohockey matches. One goroutine per match at 60 Hz physics, 20 Hz snapshot broadcast over WSS.

## Run

```bash
cp .env.example .env
make run
```

## HTTP
- `POST /internal/matches` — create a match (caller: bot service)
- `GET  /internal/matches/{id}` — match state
- `GET  /healthz`
- `WSS  /ws?matchId={id}` — gameplay socket (frontend)

See `docs/protocol.md` for the WS message protocol.

## Environment

| Var | Purpose |
|---|---|
| `ACCOUNT_BASE_URL` | Base URL of the account-management service |
| `ACCOUNT_SERVICE_TOKEN` | Service-to-service bearer token |
| `BOT_BASE_URL` | Base URL of the bot notification service |
| `BOT_SERVICE_TOKEN` | Service-to-service bearer token |
| `SERVICE_TOKEN` | Token this service accepts on `/internal/*` routes |
| `HTTP_ADDR` | Listen address (default `:8081`) |
| `FORFEIT_GRACE` | Disconnect grace period before forfeit (default `15s`) |
| `MATCH_JOIN_DEADLINE` | Time both players have to connect after match creation |
| `MATCH_DURATION` | Max match length |
| `GOAL_CAP` | Goals needed to win |

## Concurrency guarantees

- Each match runs in one goroutine at 60 Hz. Physics state is guarded by `sync.RWMutex` inside `Match`.
- Settlement is atomic: `trySettle()` / `trySetLive()` make phase transitions check-and-set under a single lock acquisition — no TOCTOU window between reading and writing phase.
- The `starter.runners` map is protected by `sync.RWMutex`; HTTP handlers (write) and WS join callbacks (read) access it concurrently.
- Disconnect detection runs during both `COUNTDOWN` and `LIVE` phases so a player leaving during the countdown correctly forfeits before the match goes live.
