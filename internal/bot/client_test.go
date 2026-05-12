package bot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotifyResult(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/matches/m1/result" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		var body MatchResult
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Reason != "score" {
			t.Fatalf("body: %+v", body)
		}
		w.WriteHeader(204)
	}))
	defer s.Close()
	c := New(s.URL, "svc")
	if err := c.NotifyResult(context.Background(), "m1", MatchResult{Reason: "score"}); err != nil {
		t.Fatal(err)
	}
}
