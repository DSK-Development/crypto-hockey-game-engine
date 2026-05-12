package ws

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/dskdev/crypto-hockey-game-engine/internal/match"
	"github.com/dskdev/crypto-hockey-game-engine/internal/protocol"
)

func buildMatch() *match.Match {
	return match.New(match.Spec{
		ID: "m1",
		Players: [2]match.Player{
			{UserID: "u1", TelegramID: 1, Username: "a"},
			{UserID: "u2", TelegramID: 2, Username: "b"},
		},
		JoinDeadline: time.Second, Duration: time.Second, GoalCap: 5,
		Now: time.Now(),
	})
}

func TestAuth_AssignsSlotForKnownUser(t *testing.T) {
	reg := match.NewRegistry()
	_ = reg.Add(buildMatch())
	fa := FakeAuth{Result: func(_ string) (string, int64, string, error) {
		return "u2", 2, "b", nil
	}}
	srv := NewServer(reg, fa, nil)
	ts := httptest.NewServer(http.HandlerFunc(srv.Handle))
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1) + "?matchId=m1"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	mustWrite(t, ctx, c, `{"type":"AUTH","auth":{"initData":"x"}}`)
	got := mustRead(t, ctx, c)
	if !strings.Contains(got, `"type":"AUTH_OK"`) || !strings.Contains(got, `"playerSlot":"B"`) {
		t.Fatalf("response: %s", got)
	}
}

func TestAuth_RejectsNonParticipant(t *testing.T) {
	reg := match.NewRegistry()
	_ = reg.Add(buildMatch())
	fa := FakeAuth{Result: func(_ string) (string, int64, string, error) {
		return "intruder", 99, "x", nil
	}}
	srv := NewServer(reg, fa, nil)
	ts := httptest.NewServer(http.HandlerFunc(srv.Handle))
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1) + "?matchId=m1"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(websocket.StatusNormalClosure, "")
	mustWrite(t, ctx, c, `{"type":"AUTH","auth":{"initData":"x"}}`)
	got := mustRead(t, ctx, c)
	if !strings.Contains(got, `"type":"AUTH_FAIL"`) {
		t.Fatalf("expected fail, got %s", got)
	}
}

func TestAuth_FailsWhenAccountErrors(t *testing.T) {
	reg := match.NewRegistry()
	_ = reg.Add(buildMatch())
	fa := FakeAuth{Result: func(_ string) (string, int64, string, error) {
		return "", 0, "", errors.New("boom")
	}}
	srv := NewServer(reg, fa, nil)
	ts := httptest.NewServer(http.HandlerFunc(srv.Handle))
	defer ts.Close()
	url := strings.Replace(ts.URL, "http", "ws", 1) + "?matchId=m1"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(websocket.StatusNormalClosure, "")
	mustWrite(t, ctx, c, `{"type":"AUTH","auth":{"initData":"x"}}`)
	got := mustRead(t, ctx, c)
	if !strings.Contains(got, `"AUTH_FAIL"`) {
		t.Fatalf("expected fail, got %s", got)
	}
	_ = protocol.TypeAuth
}

func mustWrite(t *testing.T, ctx context.Context, c *websocket.Conn, payload string) {
	t.Helper()
	if err := c.Write(ctx, websocket.MessageText, []byte(payload)); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, ctx context.Context, c *websocket.Conn) string {
	t.Helper()
	_, b, err := c.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
