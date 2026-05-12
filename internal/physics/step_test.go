package physics

import (
	"math"
	"testing"
)

const dt = 1.0 / 60.0

func TestStep_PuckMovesByVelocity(t *testing.T) {
	r := DefaultRink()
	p := &Body{Pos: Vec2{400, 200}, Vel: Vec2{60, 0}, R: 14}
	StepPuck(p, r, dt)
	if math.Abs(p.Pos.X-401) > 0.01 {
		t.Fatalf("pos: %+v", p.Pos)
	}
}

func TestStep_PuckFrictionDecays(t *testing.T) {
	r := DefaultRink()
	p := &Body{Pos: Vec2{400, 200}, Vel: Vec2{100, 0}, R: 14}
	StepPuck(p, r, dt)
	if p.Vel.X >= 100 {
		t.Fatalf("expected friction decay, got %v", p.Vel.X)
	}
}

func TestStep_PuckBouncesTopWall(t *testing.T) {
	r := DefaultRink()
	p := &Body{Pos: Vec2{400, 14.5}, Vel: Vec2{0, -100}, R: 14}
	StepPuck(p, r, dt)
	if p.Vel.Y <= 0 {
		t.Fatalf("vel.y should flip positive after bounce: %v", p.Vel.Y)
	}
	if p.Pos.Y < r.Top+p.R-1e-6 {
		t.Fatalf("pos.y should be pushed out: %v", p.Pos.Y)
	}
}

func TestStep_PuckPassesThroughGoalMouth(t *testing.T) {
	r := DefaultRink()
	// inside goal mouth, moving left past the left wall
	p := &Body{Pos: Vec2{r.Left + 5, 200}, Vel: Vec2{-200, 0}, R: 14}
	StepPuck(p, r, dt)
	// should NOT bounce; X should decrease past left wall
	if p.Vel.X > 0 {
		t.Fatalf("expected no bounce, vel: %v", p.Vel.X)
	}
}
