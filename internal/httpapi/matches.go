package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/dskdev/crypto-hockey-game-engine/internal/match"
	"github.com/go-chi/chi/v5"
)

type createMatchReq struct {
	MatchID    string `json:"matchId"`
	StakeStars int64  `json:"stakeStars"`
	Players    []struct {
		UserID     string `json:"userId"`
		TelegramID int64  `json:"telegramId"`
		Username   string `json:"username"`
	} `json:"players"`
}

func createMatch(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createMatchReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		if req.MatchID == "" || len(req.Players) != 2 {
			http.Error(w, "matchId + 2 players required", 400)
			return
		}
		spec := match.Spec{
			ID: req.MatchID, Stake: req.StakeStars,
			Players: [2]match.Player{
				{UserID: req.Players[0].UserID, TelegramID: req.Players[0].TelegramID, Username: req.Players[0].Username},
				{UserID: req.Players[1].UserID, TelegramID: req.Players[1].TelegramID, Username: req.Players[1].Username},
			},
			JoinDeadline: d.JoinDeadline,
			Duration:     d.MatchDuration,
			GoalCap:      d.GoalCap,
			Now:          d.Now(),
		}
		m := match.New(spec)
		if err := d.Registry.Add(m); err != nil {
			http.Error(w, "duplicate matchId", 409)
			return
		}
		if d.Starter != nil {
			d.Starter.Start(m)
		}
		w.WriteHeader(201)
		_ = json.NewEncoder(w).Encode(map[string]any{"matchId": req.MatchID})
	}
}

func getMatch(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		m, ok := d.Registry.Get(id)
		if !ok {
			http.Error(w, "not found", 404)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"matchId": m.ID(),
			"phase":   m.Phase(),
			"score":   m.Score(),
		})
	}
}
