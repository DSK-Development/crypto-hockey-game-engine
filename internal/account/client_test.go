package account

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthTelegram_HappyPath(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/telegram" || r.Method != http.MethodPost {
			t.Fatalf("bad route: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["initData"] != "data" {
			t.Fatalf("body: %+v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"accessToken": "jwt",
			"player":      map[string]any{"id": "u1", "telegramId": 42, "username": "alice"},
		})
	}))
	defer s.Close()
	c := New(s.URL, "svc-token")
	res, err := c.AuthTelegram(context.Background(), "data")
	if err != nil {
		t.Fatal(err)
	}
	if res.Player.ID != "u1" || res.Player.Username != "alice" {
		t.Fatalf("result: %+v", res)
	}
}

func TestPostMatchResult_RetryableOn5xx(t *testing.T) {
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			http.Error(w, "boom", 500)
			return
		}
		w.WriteHeader(204)
	}))
	defer s.Close()
	c := New(s.URL, "svc-token")
	err := c.PostMatchResult(context.Background(), MatchResult{MatchID: "m1"})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}
