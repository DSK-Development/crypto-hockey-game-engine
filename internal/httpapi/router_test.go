package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dskdev/crypto-hockey-game-engine/internal/match"
)

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	reg := match.NewRegistry()
	r := New(Deps{Registry: reg, ServiceToken: "svc"})
	return httptest.NewServer(r)
}

func TestHealthz(t *testing.T) {
	s := newServer(t)
	defer s.Close()
	res, err := http.Get(s.URL + "/healthz")
	if err != nil || res.StatusCode != 200 {
		t.Fatalf("err=%v status=%d", err, res.StatusCode)
	}
}

func TestCreateMatch_RequiresServiceToken(t *testing.T) {
	s := newServer(t)
	defer s.Close()
	body, _ := json.Marshal(map[string]any{
		"matchId": "m1", "stakeStars": 10,
		"players": []map[string]any{
			{"userId": "u1", "telegramId": 1, "username": "a"},
			{"userId": "u2", "telegramId": 2, "username": "b"},
		},
	})
	res, _ := http.Post(s.URL+"/internal/matches", "application/json", bytes.NewReader(body))
	if res.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestCreateMatch_CreatesAndGet(t *testing.T) {
	s := newServer(t)
	defer s.Close()
	body, _ := json.Marshal(map[string]any{
		"matchId": "m1", "stakeStars": 10,
		"players": []map[string]any{
			{"userId": "u1", "telegramId": 1, "username": "a"},
			{"userId": "u2", "telegramId": 2, "username": "b"},
		},
	})
	req, _ := http.NewRequest("POST", s.URL+"/internal/matches", bytes.NewReader(body))
	req.Header.Set("X-Service-Token", "svc")
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 201 {
		t.Fatalf("create: %v %d", err, res.StatusCode)
	}
	// GET
	req, _ = http.NewRequest("GET", s.URL+"/internal/matches/m1", nil)
	req.Header.Set("X-Service-Token", "svc")
	res, _ = http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		t.Fatalf("get: %d", res.StatusCode)
	}
}
