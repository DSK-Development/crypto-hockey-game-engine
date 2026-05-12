package httpapi

import (
	"net/http"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/match"
	"github.com/go-chi/chi/v5"
)

type MatchStarter interface {
	Start(m *match.Match)
}

type Deps struct {
	Registry          *match.Registry
	ServiceToken      string
	Starter           MatchStarter
	JoinDeadline      time.Duration
	MatchDuration     time.Duration
	GoalCap           int
	Now               func() time.Time
}

func New(d Deps) http.Handler {
	if d.Now == nil {
		d.Now = time.Now
	}
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })
	r.Route("/internal", func(r chi.Router) {
		r.Use(requireServiceToken(d.ServiceToken))
		r.Post("/matches", createMatch(d))
		r.Get("/matches/{id}", getMatch(d))
	})
	return r
}

func requireServiceToken(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Service-Token") != token {
				http.Error(w, "unauthorized", 401)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
