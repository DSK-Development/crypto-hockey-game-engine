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
