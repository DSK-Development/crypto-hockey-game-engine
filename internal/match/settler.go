package match

import (
	"context"

	"github.com/dskdev/crypto-hockey-game-engine/internal/account"
	"github.com/dskdev/crypto-hockey-game-engine/internal/bot"
)

type accountPoster interface {
	PostMatchResult(ctx context.Context, r account.MatchResult) error
}

type botNotifier interface {
	NotifyResult(ctx context.Context, matchID string, r bot.MatchResult) error
}

type defaultSettler struct {
	acc accountPoster
	bot botNotifier
}

func NewSettler(a accountPoster, b botNotifier) *defaultSettler { return &defaultSettler{acc: a, bot: b} }

func (s *defaultSettler) Settle(ctx context.Context, m *Match, winnerUserID, reason string) {
	if m.Phase() == PhaseSettled {
		return
	}
	m.SetPhase(PhaseSettled)
	score := m.Score()
	parts := make([]account.Participant, 0, 2)
	for i, p := range m.Spec().Players {
		placement := 2
		if p.UserID == winnerUserID {
			placement = 1
		}
		// for ties on timeout, both placements = 1 (handled by callers if needed); v1 leaves 2/2 if no winner.
		_ = i
		parts = append(parts, account.Participant{UserID: p.UserID, Placement: placement})
	}
	_ = s.acc.PostMatchResult(ctx, account.MatchResult{
		MatchID: m.ID(), Stake: m.Spec().Stake,
		Participants: parts, Reason: reason,
	})
	var winnerPtr *string
	if winnerUserID != "" {
		w := winnerUserID
		winnerPtr = &w
	}
	_ = s.bot.NotifyResult(ctx, m.ID(), bot.MatchResult{
		WinnerUserID: winnerPtr, Reason: reason, ScoreA: score.A, ScoreB: score.B,
	})
}
