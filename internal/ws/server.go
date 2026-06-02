package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/dskdev/crypto-hockey-game-engine/internal/match"
	"github.com/dskdev/crypto-hockey-game-engine/internal/protocol"
)

type AuthResolver interface {
	AuthTelegram(ctx context.Context, initData string) (userID string, telegramID int64, username string, err error)
}

type Server struct {
	reg      *match.Registry
	auth     AuthResolver
	onJoin   func(m *match.Match, slot int, conn *Conn)
	mu       sync.Mutex
}

func NewServer(reg *match.Registry, auth AuthResolver, onJoin func(*match.Match, int, *Conn)) *Server {
	return &Server{reg: reg, auth: auth, onJoin: onJoin}
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	matchID := r.URL.Query().Get("matchId")
	m, ok := s.reg.Get(matchID)
	if !ok {
		http.Error(w, "unknown match", 404)
		return
	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}
	conn := newConn(c)
	defer conn.Close()

	// Read first message: must be AUTH.
	ctx := r.Context()
	var env protocol.ClientEnvelope
	if err := conn.ReadJSON(ctx, &env); err != nil {
		_ = conn.WriteJSON(ctx, protocol.ServerEnvelope{Type: protocol.TypeAuthFail, AuthFail: &protocol.AuthFailPayload{Reason: "bad message"}})
		return
	}
	if env.Type != protocol.TypeAuth || env.Auth == nil {
		_ = conn.WriteJSON(ctx, protocol.ServerEnvelope{Type: protocol.TypeAuthFail, AuthFail: &protocol.AuthFailPayload{Reason: "first message must be AUTH"}})
		return
	}
	userID, telegramID, username, err := s.auth.AuthTelegram(ctx, env.Auth.InitData)
	if err != nil {
		slog.Error("ws auth failed", "err", err, "initDataLen", len(env.Auth.InitData))
		_ = conn.WriteJSON(ctx, protocol.ServerEnvelope{Type: protocol.TypeAuthFail, AuthFail: &protocol.AuthFailPayload{Reason: "auth failed: " + err.Error()}})
		return
	}
	slot, opp, err := assignSlot(m, userID, telegramID)
	if err != nil {
		slog.Error("ws assignSlot failed", "err", err, "userID", userID, "telegramID", telegramID)
		_ = conn.WriteJSON(ctx, protocol.ServerEnvelope{Type: protocol.TypeAuthFail, AuthFail: &protocol.AuthFailPayload{Reason: err.Error()}})
		return
	}
	conn.userID, conn.telegramID, conn.username = userID, telegramID, username
	conn.slot = slot
	_ = conn.WriteJSON(ctx, protocol.ServerEnvelope{
		Type:   protocol.TypeAuthOK,
		AuthOK: &protocol.AuthOKPayload{PlayerSlot: slotLabel(slot), Opponent: opp},
	})
	if s.onJoin != nil {
		s.onJoin(m, slot, conn)
	}
	// Hand off to onJoin which owns subsequent reads. Block until conn closes.
	conn.WaitClosed()
}

func assignSlot(m *match.Match, userID string, telegramID int64) (int, protocol.Opponent, error) {
	players := m.Spec().Players
	for i, p := range players {
		if p.UserID == userID || p.TelegramID == telegramID {
			opp := players[1-i]
			return i, protocol.Opponent{Username: opp.Username, TelegramID: opp.TelegramID}, nil
		}
	}
	return 0, protocol.Opponent{}, errors.New("not a participant")
}

func slotLabel(i int) string {
	if i == 0 {
		return "A"
	}
	return "B"
}

// FakeAuth helper for tests.
type FakeAuth struct{ Result func(string) (string, int64, string, error) }

func (f FakeAuth) AuthTelegram(_ context.Context, initData string) (string, int64, string, error) {
	return f.Result(initData)
}

var _ = json.Marshal // keep encoding/json referenced if compiler complains during partial implementations
