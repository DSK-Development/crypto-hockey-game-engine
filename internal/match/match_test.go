package match

import (
	"testing"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/physics"
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

func TestTick_GoalCap_EndedWhenCapReached(t *testing.T) {
	m := New(testSpec(t, 3)) // cap at 3
	m.SetPhase(PhaseLive)

	goals := 0
	for goals < 3 {
		// Place puck near right goal mouth (A scores on right side)
		m.SetPuck(physics.Body{Pos: physics.Vec2{X: 795, Y: 200}, Vel: physics.Vec2{X: 300, Y: 0}, R: 14})
		for i := 0; i < 60; i++ {
			if scorer := m.Tick(1.0 / 60.0); scorer != "" {
				goals++
				break
			}
		}
	}

	s := m.Score()
	if s.A != 3 {
		t.Fatalf("expected A to score 3, got %+v", s)
	}
	capReached := s.A >= 3 || s.B >= 3
	if !capReached {
		t.Fatal("goal cap not reached")
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
