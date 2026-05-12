package match

import (
	"testing"
	"time"
)

func TestNewMatch_Pending(t *testing.T) {
	m := New(Spec{
		ID: "m1", Stake: 50,
		Players: [2]Player{
			{UserID: "u1", TelegramID: 1, Username: "a"},
			{UserID: "u2", TelegramID: 2, Username: "b"},
		},
		JoinDeadline: 60 * time.Second,
		Duration:     120 * time.Second,
		GoalCap:      5,
		Now:          time.Unix(1000, 0),
	})
	if m.Phase() != PhasePending {
		t.Fatalf("phase: %s", m.Phase())
	}
	if m.JoinDeadlineAt() != time.Unix(1060, 0) {
		t.Fatal("join deadline")
	}
}

func TestApplyGoal_IncrementsAndDetectsWin(t *testing.T) {
	m := New(testSpec(t, 5))
	m.SetPhase(PhaseLive)
	for i := 0; i < 4; i++ {
		if win := m.ApplyGoal("A"); win {
			t.Fatalf("premature win after %d", i+1)
		}
	}
	if !m.ApplyGoal("A") {
		t.Fatal("expected win on 5th goal")
	}
	if s := m.Score(); s.A != 5 || s.B != 0 {
		t.Fatalf("score: %+v", s)
	}
}

func testSpec(t *testing.T, cap int) Spec {
	t.Helper()
	return Spec{
		ID: "m1", Stake: 10,
		Players: [2]Player{
			{UserID: "u1", TelegramID: 1, Username: "a"},
			{UserID: "u2", TelegramID: 2, Username: "b"},
		},
		JoinDeadline: time.Second, Duration: time.Second,
		GoalCap: cap, Now: time.Unix(1000, 0),
	}
}
