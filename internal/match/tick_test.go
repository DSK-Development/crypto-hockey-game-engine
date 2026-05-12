package match

import (
	"testing"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/physics"
)

func TestTick_AdvancesPuck(t *testing.T) {
	m := New(testSpec(t, 5))
	m.SetPhase(PhaseLive)
	m.SetPuck(physics.Body{Pos: physics.Vec2{X:400, Y:200}, Vel: physics.Vec2{X:120, Y:0}, R: 14})
	m.Tick(1.0 / 60.0)
	snap := m.Snapshot(time.Unix(0, 0))
	if snap.Puck.X <= 400 {
		t.Fatalf("puck did not advance: %+v", snap.Puck)
	}
}

func TestTick_DetectsGoalAndResetsPuck(t *testing.T) {
	m := New(testSpec(t, 5))
	m.SetPhase(PhaseLive)
	// puck on right edge, moving right inside goal mouth
	m.SetPuck(physics.Body{Pos: physics.Vec2{X:795, Y:200}, Vel: physics.Vec2{X:200, Y:0}, R: 14})
	for i := 0; i < 30 && m.Score().A == 0; i++ {
		m.Tick(1.0 / 60.0)
	}
	if m.Score().A != 1 {
		t.Fatalf("expected A goal, got %+v", m.Score())
	}
}
